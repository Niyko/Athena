package main

import (
	"os"

	"github.com/fatih/color"
)

func main() {
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
	color.New(color.FgCyan, color.Bold).Printf("🪅  Athena • v1.02\n")
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
	color.New(color.FgGreen).Printf("\trecreate-sqlite")
	color.White(" - Recreate the SQlite database of Athena\n")
	color.New(color.FgGreen).Printf("\tuninstall")
	color.White(" - Uninstall the things that needed for Athena\n")
	color.New(color.FgGreen).Printf("\thelp")
	color.White(" - Display help details of all commands")
}
