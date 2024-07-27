package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/fatih/color"
)

func runClickhouseMigration() {
	config := getConfig()

	if !config.ClickHouse {
		color.Yellow("⚠️  Clickhouse log is disabled in config")
		return
	}

	clickhouseDB := getClickhouseConnection()

	if checkClickhouseTableExists(clickhouseDB, config.ClickHouseDatabase, config.ClickHouseTableName) {
		color.Yellow("✅  Clickhouse log table already exists")
	} else {
		createTableQuery := fmt.Sprintf(`
			CREATE TABLE %s.%s (
				"count" Int64,
				"table" String,
				"timestamp" DateTime
			)
			ENGINE = MergeTree
			ORDER BY timestamp
			TTL timestamp + toIntervalHour(%d);
		`, config.ClickHouseDatabase, config.ClickHouseTableName, config.ClickHouseTableTTL)

		_, error := clickhouseDB.Exec(createTableQuery)
		if error != nil {
			color.Red("Error while creating table in clickhouse (%s)", error)
			sentry.CaptureException(error)
			os.Exit(0)
		}

		color.Green("✅  Clickhouse log table has been created")
	}

	defer clickhouseDB.Close()

	color.White("  ")
	color.White("✅  Clickhouse migration has been completed")
}

func checkClickhouseTableExists(db *sql.DB, databaseName, tableName string) (bool) {
    var name string
    query := "SELECT name FROM system.tables WHERE database = ? AND name = ?"
    error := db.QueryRow(query, databaseName, tableName).Scan(&name)

    if error != nil {
        if error == sql.ErrNoRows {
            return false
        }
        return false
    }
    return true
}

func getClickhouseConnection() *sql.DB {
	config := getConfig()

	clickhouseDB := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{config.ClickHouseHost},
		Auth: clickhouse.Auth{
			Database: config.ClickHouseDatabase,
			Username: config.ClickHouseUsername,
			Password: config.ClickHousePassword,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 30 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Protocol:  clickhouse.HTTP,
	})

	if error := clickhouseDB.Ping(); error != nil {
        color.Red("Error while connecting to clickhouse (%s)", error)
		sentry.CaptureException(error)
		os.Exit(0)
    }

	return clickhouseDB
}