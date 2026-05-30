package main

import (
	"os"
	"context"
	"time"
	"crypto/tls"

	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func getKafkaWriter() *kafka.Writer {
	config := getConfig()

	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
		DualStack: true,
	}

	if config.KafkaEnableTLS {
		dialer.TLS = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	if config.KafkaSASLMechanisms != "NONE" {
		var mechanism sasl.Mechanism
		var error error
		
		switch config.KafkaSASLMechanisms {
			case "SASL-SCRAM-SHA-256":
				mechanism, error = scram.Mechanism(
					scram.SHA256,
					config.KafkaSASLUsername,
					config.KafkaSASLPassword,
				)

				if error != nil {
					color.Red("Error while initializing SASL-SCRAM-SHA-256. (%s)", error)
					os.Exit(0)
				}

			case "SASL-SCRAM-SHA-512":
				mechanism, error = scram.Mechanism(
					scram.SHA256,
					config.KafkaSASLUsername,
					config.KafkaSASLPassword,
				)

				if error != nil {
					color.Red("Error while initializing SASL-SCRAM-SHA-512. (%s)", error)
					os.Exit(0)
				}

			default:
				mechanism = plain.Mechanism{
					Username: config.KafkaSASLUsername,
					Password: config.KafkaSASLPassword,
				}
		}

		dialer.SASLMechanism = mechanism
	}

	kafkaConfig := kafka.WriterConfig{
		Brokers: []string{config.KafkaHost},
		Topic: config.KafkaTopic,
		BatchTimeout: 50 * time.Millisecond,
		Dialer: dialer,
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