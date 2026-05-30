package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/segmentio/kafka-go"
	"github.com/yosuke-furukawa/json5/encoding/json5"
	"github.com/fatih/color"
)

type Config struct {
	DBHost         string `json:"dbHost"`
	DBPort         int    `json:"dbPort"`
	DBUser         string `json:"dbUser"`
	DBPassword     string `json:"dbPassword"`
	DBName         string `json:"dbName"`
	KafkaHost      string `json:"kafkaHost"`
	KafkaTopic     string `json:"kafkaTopic"`
	PollInterval   int    `json:"pollInterval"`
	FetchLimit     int    `json:"fetchLimit"`
	SkippedTables  []string `json:"skippedTables"`
}

func TestIntegration(t *testing.T) {
	config := loadConfig(t)

	setupCmd := exec.Command("go", "run", ".", "setup")
	setupCmd.Dir = filepath.Join(filepath.Dir(getTestPath()), "..")
	setupCmd.Env = append(os.Environ(), "GORUN=1")

	if err := setupCmd.Run(); err != nil {
		t.Fatal(color.RedString("⛔️ Failed to run setup: %v", err))
	} else {
		t.Log(color.GreenString("✅ Successfully run the Athena setup"))
	}

	cmd := exec.Command("go", "run", ".", "run")
	cmd.Dir = filepath.Join(filepath.Dir(getTestPath()), "..")
	cmd.Env = append(os.Environ(), "GORUN=1")

	if err := cmd.Start(); err != nil {
		t.Fatalf(color.RedString("⛔️ Failed to start Athena daemon: %v", err))
	} else{
		t.Log(color.GreenString("✅ Successfully started the Athena daemon"))
	}

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	time.Sleep(3 * time.Second)

	testProductID := rand.Intn(90000) + 10000
	testProductName := "Test product"

	db := connectToMSSQL(t, config)
	defer db.Close()

	insertTestProduct(t, db, testProductID, testProductName)

	time.Sleep(time.Duration(config.PollInterval+2) * time.Second)

	verifyKafkaMessage(t, config, testProductID)

	deleteTestProduct(t, db, testProductID)
}

func loadConfig(t *testing.T) Config {
	configPath := filepath.Join(filepath.Dir(getTestPath()), "..", "config.json")

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(color.RedString("⛔️ Failed to read config file: %v", err))
	}

	var config Config
	if err := json5.Unmarshal(content, &config); err != nil {
		t.Fatal(color.RedString("⛔️ Failed to parse config file: %v", err))
	}

	return config
}

func connectToMSSQL(t *testing.T, config Config) *sql.DB {
	connectionString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		t.Fatal(color.RedString("⛔️ Failed to open MSSQL connection (%v)", err))
	}

	if err := db.Ping(); err != nil {
		t.Fatal(color.RedString("⛔️ Failed to ping MSSQL (%v)", err))
	}

	return db
}

func insertTestProduct(t *testing.T, db *sql.DB, productID int, productName string) {
	query := fmt.Sprintf(
		"INSERT INTO products (id, name) VALUES (%d, '%s')",
		productID,
		productName,
	)

	if _, err := db.Exec(query); err != nil {
		t.Fatalf(color.RedString("⛔️ Failed to insert test product (%v)", err))
	}

	t.Log(color.GreenString("✅ Inserted test product ID (%d) and name (%s)", productID, productName))
}

func verifyKafkaMessage(t *testing.T, config Config, expectedProductID int) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{config.KafkaHost},
		Topic:          config.KafkaTopic,
		CommitInterval: time.Second,
	})
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				t.Fatal(color.RedString("⛔️ Timeout: Did not receive Kafka message for product ID (%d) within 20 seconds", expectedProductID))
			}
			t.Log(color.RedString("⛔️ Error reading message (%v)", err))
			continue
		}

		var kafkaData map[string]interface{}
		if err := json.Unmarshal(msg.Value, &kafkaData); err != nil {
			t.Log(color.RedString("⛔️ Failed to unmarshal message (%v)", err))
			continue
		}

		tableName, ok := kafkaData["tableName"].(string)
		if !ok || tableName != "products" {
			continue
		}

		rowID, ok := kafkaData["rowId"]
		if !ok {
			continue
		}

		idVal, ok := rowID.(float64)
		if !ok {
			continue
		}

		if int(idVal) == expectedProductID {
			t.Log(color.GreenString("✅ Successfully verified Kafka message for product ID (%d)", expectedProductID))
			return
		}
	}
}

func deleteTestProduct(t *testing.T, db *sql.DB, productID int) {
	query := fmt.Sprintf("DELETE FROM products WHERE id = %d", productID)
	if _, err := db.Exec(query); err != nil {
		t.Log(color.YellowString("⚠️ Warning: Failed to delete test product (%v)", err))
	}
}

func getTestPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("⛔️  Failed to get current file path")
	}
	return file
}
