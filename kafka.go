package main

import (
	"context"
	"time"
	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"crypto/tls"
)

func getKafkaWriter() *kafka.Writer {
	config := getConfig()
	mechanism := plain.Mechanism{
		Username: config.KafkaSASLUsername,
		Password: config.KafkaSASLPassword,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	kafkaConfig := kafka.WriterConfig{
		Brokers: []string{config.KafkaHost},
		Topic: config.KafkaTopic,
		BatchTimeout: 50 * time.Millisecond,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			Timeout: 10 * time.Second,
			DualStack: true,
			TLS: tlsConfig,
		},
	}

	kafkaWriter := kafka.NewWriter(kafkaConfig)

	return kafkaWriter
}

func sendMessageToKafka(message kafka.Message, kafkaWriter *kafka.Writer) {
	error := kafkaWriter.WriteMessages(context.Background(), message)
	if error != nil {
		color.Red("Error while connecting to Kafka (%s)", error)
	}
}