package main

import (
	"os"

  	"github.com/getsentry/sentry-go"
	"github.com/fatih/color"
)

func main() {
	config := getConfig()

	sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRYDNS"),
		TracesSampleRate: 1.0,
		Transport: sentry.NewHTTPSyncTransport(),
	})

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("character", map[string]interface{}{
			"uuid": config.Uuid,
			"argument": os.Args,
		})
	})

	sentry.CaptureMessage("athena has been started")

	if len(os.Args) > 1 {
		if os.Args[1] == "run" {
			startDemon()
		} else if os.Args[1] == "setup" {
			installAthena()
		} else if os.Args[1] == "uninstall" {
			uninstallAthena()
		} else if os.Args[1] == "add-cdc" {
			setupMSSQLCDC()
		} else if os.Args[1] == "remove-cdc" {
			disableMSSQLDatabaseCDC()
		} else if os.Args[1] == "clear-cdc-history" {
			removeCDCHistory()
		} else if os.Args[1] == "recreate-clickhouse" {
			runClickhouseMigration()
		} else if os.Args[1] == "recreate-sqlite" {
			runSqliteMigration()
		} else if os.Args[1] == "help" {
			printHelp()
		} else {
			color.Red("⛔️  Command is not found")
			color.White("  ")
			printHelp()
		}
	} else {
		printHelp()
	}
}

func installAthena() {
	setupMSSQLCDC()
	runClickhouseMigration()
	runSqliteMigration()

	color.White("  ")
	color.Cyan("🪅  Athena is ready to be started")
}

func uninstallAthena() {
	disableMSSQLDatabaseCDC()
	deleteSqliteDatabase()

	color.White("  ")
	color.Cyan("🪅 Athena is uninstalled successfully")
}

func printHelp() {
	color.New(color.FgCyan, color.Bold).Printf("🪅  Athena • v1.04\n")
	color.White("Go to https://github.com/cristalhq/acmd for more info")
	color.Yellow("\n Usage:")
	color.New(color.FgGreen).Printf("\trun")
	color.White(" - To run the change data capture demon\n")
	color.New(color.FgGreen).Printf("\tsetup")
	color.White(" - Setup all required things for Athena\n")
	color.New(color.FgGreen).Printf("\tadd-cdc")
	color.White(" - Setup change data caputure in MSSQL database\n")
	color.New(color.FgGreen).Printf("\tremove-cdc")
	color.White(" - Remove change data caputure in MSSQL database\n")
	color.New(color.FgGreen).Printf("\tclear-cdc-history")
	color.White(" - Clear existing CDC history of database\n")
	color.New(color.FgGreen).Printf("\trecreate-sqlite")
	color.White(" - Recreate the SQlite database of Athena\n")
	color.New(color.FgGreen).Printf("\trecreate-clickhouse")
	color.White(" - Recreate the Clickhouse log tables\n")
	color.New(color.FgGreen).Printf("\tuninstall")
	color.White(" - Uninstall the things that needed for Athena\n")
	color.New(color.FgGreen).Printf("\thelp")
	color.White(" - Display help details of all commands")
}
