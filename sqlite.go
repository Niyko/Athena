package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CDCLastTableLSN struct {
	gorm.Model
	TableName string `gorm:"unique"`
	LSN       string
}

func runSqliteMigration() {
	sqliteFileName := "db.sqlite"

	deleteSqliteDatabase()

	sqliteFile, error := os.Create(sqliteFileName)
	if error != nil {
		color.Red("Error while creating SQlite file (%s)", error)
		os.Exit(0)
	}
	defer sqliteFile.Close()

	color.Green("✅  SQlite file has been created")

	sqliteDatabase := getSqliteConnection()

	sqliteDatabase.Migrator().CreateTable(&CDCLastTableLSN{})

	color.White("  ")
	color.White("✅  SQlite migration has been completed")
}

func deleteSqliteDatabase() {
	sqliteFileName := "db.sqlite"
	if _, error := os.Stat(sqliteFileName); error == nil {
		error := os.Remove(sqliteFileName)
		if error != nil {
			color.Red("Error while deleting SQlite file (%s)", error)
			os.Exit(0)
		}

		color.Yellow("✅  SQlite file has been deleted")
	}
}

func getSqliteConnection() *gorm.DB {
	database, error := gorm.Open(sqlite.Open(getExePath() + "db.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if error != nil {
		color.Red("Error while connecting to SQLite (%s)", error)
		os.Exit(0)
	}

	return database
}
