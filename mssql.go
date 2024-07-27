package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/getsentry/sentry-go"
	"github.com/fatih/color"
)

func setupMSSQLCDC() {
	config := getConfig()
	database := getMSSQLConnection()
	defer database.Close()

	if checkMSSQLDatabaseCDCEnabled(config, database) {
		color.Yellow("✅  CDC is already enabled in database level MSSQL")
	}

	if !checkMSSQLDatabaseCDCEnabled(config, database) {
		result, error := database.Query(fmt.Sprintf(`
			sp_changedbowner '%s';
			EXEC sys.sp_cdc_enable_db;
		`, config.DBUser))
		_ = result

		if error != nil {
			color.Red("Error while checking CDC enabled in MSSQL (%s)", error)
			sentry.CaptureException(error)
			os.Exit(0)
		}
		defer result.Close()

		color.Green("✅  CDC is enabled in database level MSSQL\n")
	}

	if !checkMSSQLDatabaseCDCEnabled(config, database) {
		color.Red("Error while enabling CDC in MSSQL")
		os.Exit(0)
	}

	color.White("  ")

	tablesList := getAllTablesInMSSQL(database)

	for _, tableName := range tablesList {
		if contains(config.SkippedTables, tableName) {
			color.Cyan("🟦  Table [%s] is skipped from tacking as added in config file", tableName)
			continue
		}

		if checkMSSQLTableCDCEnabled(tableName, database) {
			color.Yellow("✅  CDC is already enabled in [%s] table level MSSQL", tableName)
		}

		if !checkMSSQLTableCDCEnabled(tableName, database) {
			result, error := database.Query(fmt.Sprintf(`
				EXEC sys.sp_cdc_enable_table 
					@source_schema = N'dbo', 
					@source_name   = N'%s', 
					@role_name     = NULL, 
					@supports_net_changes = 1
			`, tableName))
			_ = result

			if error != nil {
				color.Red("Error while enabling CDC for [%s] table in MSSQL (%s)", tableName, error)
				sentry.CaptureException(error)
				os.Exit(0)
			}
			defer result.Close()

			color.Green("✅  CDC is enabled in [%s] table level MSSQL", tableName)
		}

		if !checkMSSQLTableCDCEnabled(tableName, database) {
			color.Red("Error while enabling CDC in [%s] table level MSSQL", tableName)
			os.Exit(0)
		}
	}

	color.White("  ")
	color.White("✅  CDC installation completed")
}

func removeCDCHistory() {
	config := getConfig()
	database := getMSSQLConnection()
	defer database.Close()

	tablesList := getAllTablesInMSSQL(database)

	for _, tableName := range tablesList {
		if contains(config.SkippedTables, tableName) {
			color.Cyan("🟦  Table [%s] is skipped from tacking as added in config file", tableName)
			continue
		}

		if checkMSSQLTableCDCEnabled(tableName, database) {
			result, error := database.Query(fmt.Sprintf(`
				DELETE FROM cdc.dbo_%s_CT;
			`, tableName))
			_ = result

			if error != nil {
				color.Red("Error while clearing CDC history for [%s] table in MSSQL (%s)", tableName, error)
				sentry.CaptureException(error)
				os.Exit(0)
			}
			defer result.Close()

			color.Green("✅  CDC history is cleared in [%s] table level MSSQL", tableName)
		}

		if !checkMSSQLTableCDCEnabled(tableName, database) {
			color.Red("Error while clearing CDC history in [%s] table level MSSQL", tableName)
			os.Exit(0)
		}
	}

	color.White("  ")

	runSqliteMigration()

	color.White("  ")
	color.White("✅  CDC history has been cleared")
}

func getAllTablesInMSSQL(database *sql.DB) []string {
	rows, error := database.QueryContext(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND TABLE_SCHEMA = 'dbo'")
	if error != nil {
		color.Red("Error while getting table list in MSSQL (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var tableName string
		if error := rows.Scan(&tableName); error != nil {
			color.Red("Error while getting table list in MSSQL (%s)", error)
			sentry.CaptureException(error)
			os.Exit(0)
		}
		if tableName != "sysdiagrams" && tableName != "systranschemas" {
			tables = append(tables, tableName)
		}
	}

	if error := rows.Err(); error != nil {
		color.Red("Error while getting table list in MSSQL (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
	}

	return tables
}

func checkMSSQLTableCDCEnabled(tableName string, database *sql.DB) bool {
	rows, error := database.Query(fmt.Sprintf(`
		SELECT name, is_tracked_by_cdc
		FROM sys.tables
		WHERE name = '%s' AND is_tracked_by_cdc = 1;
	`, tableName))

	if error != nil {
		color.Red("Error while checking CDC enabled in [%s] table level in MSSQL (%s)", tableName, error)
		sentry.CaptureException(error)
		os.Exit(0)
	}
	defer rows.Close()

	return rows.Next()
}

func checkMSSQLDatabaseCDCEnabled(config Config, database *sql.DB) bool {
	rows, error := database.Query(fmt.Sprintf(`
		SELECT name, is_cdc_enabled
		FROM sys.databases
		WHERE is_cdc_enabled = 1 AND name = '%s';
	`, config.DBName))

	if error != nil {
		color.Red("Error while checking CDC enabled in database level in MSSQL (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
	}
	defer rows.Close()

	return rows.Next()
}

func disableMSSQLDatabaseCDC() {
	config := getConfig()
	mssqlDatabase := getMSSQLConnection()
	defer mssqlDatabase.Close()

	result, error := mssqlDatabase.Query(fmt.Sprintf(`
		sp_changedbowner '%s';
		EXEC sys.sp_cdc_disable_db;
	`, config.DBUser))
	_ = result

	if error != nil {
		color.Red("Error while checking CDC disabling in MSSQL (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
	}
	defer result.Close()

	color.Green("✅  CDC has been disabled in MSSQL database level")
}

func getMSSQLConnection() *sql.DB {
	config := getConfig()
	connectionString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBPort,
		config.DBName,
	)

	database, error := sql.Open("sqlserver", connectionString)
	if error != nil {
		color.Red("Error while connecting to MSSQL (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
	}

	error = database.Ping()
	if error != nil {
		color.Red("Error while pinging to MSSQL (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
	}

	return database
}
