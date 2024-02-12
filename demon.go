package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func startDemon() {
	config := getConfig()
	mssqlDatabase := getMSSQLConnection()
	sqliteDatabase := getSqliteConnection()
	kafkaWriter := getKafkaWriter()
	tables := getAllTablesInMSSQL(mssqlDatabase)

	defer mssqlDatabase.Close()
	defer kafkaWriter.Close()

	color.Green("🪅  Athena demon has started")

	for {
		for _, tableName := range tables {
			if !contains(config.SkippedTables, tableName) {
				pollChanges(tableName, mssqlDatabase, sqliteDatabase, kafkaWriter, config)
			}
		}

		time.Sleep(time.Duration(config.PollInterval) * time.Second)
	}
}

func pollChanges(tableName string, mssqlDatabase *sql.DB, sqliteDatabase *gorm.DB, kafkaWriter *kafka.Writer, config Config) {
	var cdcLastTableLSN CDCLastTableLSN
	totalCDCChangesSend := 0
	sqliteDatabaseResult := sqliteDatabase.First(&cdcLastTableLSN, "table_name = ?", tableName)
	cdcSqliteQuery := fmt.Sprintf(`SELECT TOP (%d) * FROM cdc.dbo_%s_CT`, config.FetchLimit, tableName)

	if sqliteDatabaseResult.RowsAffected == 1 {
		cdcSqliteQuery = fmt.Sprintf(`SELECT TOP (%d) * FROM cdc.dbo_%s_CT WHERE __$start_lsn>0x%s`, config.FetchLimit, tableName, cdcLastTableLSN.LSN)
	}

	rows, error := mssqlDatabase.Query(cdcSqliteQuery)

	if error != nil {
		color.Red("Error while fetching CDC changes from MSSQL. (%s) and table (%s)", error, tableName)
		os.Exit(0)
	}
	defer rows.Close()

	columns, error := rows.Columns()
	if error != nil {
		color.Red("Error while fetching columns CDC changes from MSSQL (%s)", error)
		os.Exit(0)
	}

	columnValues := make([]interface{}, len(columns))
	for i := range columns {
		columnValues[i] = new(interface{})
	}

	for rows.Next() {
		rowValues := make(map[string]interface{})
		error := rows.Scan(columnValues...)
		if error != nil {
			color.Red("Error while fetching column values CDC changes from MSSQL (%s)", error)
			os.Exit(0)
		}

		for i, column := range columns {
			columnValue := *(columnValues[i].(*interface{}))
			rowValues[column] = columnValue
		}

		operationCode := rowValues["__$operation"].(int64)

		if operationCode != 0 && operationCode != 3 && operationCode != 5 {
			kafkaData := map[string]interface{}{
				"tableName": tableName,
				"operationName": convertCDCOpertionCode(operationCode),
				"operationCode": operationCode,
				"rowId":     rowValues["id"],
				"rowValues": rowValues, 
			}

			kafkaDataString, error := json.Marshal(kafkaData)
			if error != nil {
				color.Red("Error while converting kafka data interface to json (%s)", error)
				os.Exit(0)
			}

			kafkaMessage := kafka.Message{
				Key:   []byte(hex.EncodeToString(rowValues["__$start_lsn"].([]uint8))),
				Value: []byte(string(kafkaDataString)),
			}

			sendMessageToKafka(kafkaMessage, kafkaWriter)

			sqliteDatabase.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "table_name"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"lsn": hex.EncodeToString(rowValues["__$start_lsn"].([]uint8)),
				}),
			}).Create(&[]CDCLastTableLSN{{
				TableName: tableName,
				LSN:       hex.EncodeToString(rowValues["__$start_lsn"].([]uint8)),
			}})

			totalCDCChangesSend = totalCDCChangesSend + 1
		}
	}

	if totalCDCChangesSend > 0 {
		color.White("Total %d CDC changes moved to kafka from [%s]", totalCDCChangesSend, tableName)
	}

	if error := rows.Err(); error != nil {
		color.Red("Error while iterating over rows from CDC changes from MSSQL (%s)", error)
		os.Exit(0)
	}
}

func convertCDCOpertionCode(code int64) string {
	if code == 1 {
		return "insert"
	} else if code == 2 {
		return "delete"
	} else if code == 4 {
		return "update"
	} else {
		return "INVALID"
	}
}