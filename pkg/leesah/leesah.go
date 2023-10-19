package leesah

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Rapid struct {
	writer   *kafka.Writer
	reader   *kafka.Reader
	ctx      context.Context
	teamName string
	log      *slog.Logger
}

type RapidHandlers struct {
	HandleQuestion   func(question Question) (Answer, error)
	HandleAssessment func(assessment Assessment) error
}

type rapidConfig struct {
	Brokers             string
	Topic               string
	GroupID             string
	KafkaCertPath       string
	KafkaPrivateKeyPath string
	KafkaCAPath         string
}

var config = rapidConfig{}

func init() {
	flag.StringVar(&config.Brokers, "brokers", os.Getenv("KAFKA_BROKERS"), "Kafka broker")
	flag.StringVar(&config.Topic, "topic", os.Getenv("KAFKA_TOPIC"), "Kafka topic")
	flag.StringVar(&config.GroupID, "group-id", os.Getenv("KAFKA_GROUP_ID"), "Kafka group ID")
	flag.StringVar(&config.KafkaCertPath, "kafka-cert-path", os.Getenv("KAFKA_CERTIFICATE_PATH"), "Path to Kafka certificate")
	flag.StringVar(&config.KafkaPrivateKeyPath, "kafka-private-key-path", os.Getenv("KAFKA_PRIVATE_KEY_PATH"), "Path to Kafka private key")
	flag.StringVar(&config.KafkaCAPath, "kafka-ca-path", os.Getenv("KAFKA_CA_PATH"), "Path to Kafka CA certificate")
}

func NewRapid(teamName string, logger *slog.Logger) (*Rapid, error) {
	flag.Parse()

	keypair, err := tls.LoadX509KeyPair(config.KafkaCertPath, config.KafkaPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load Access Key and/or Access Certificate: %s", err)
	}

	caCert, err := os.ReadFile(config.KafkaCAPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA Certificate file: %s", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to parse CA Certificate file: %s", err)
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{keypair},
			RootCAs:      caCertPool,
		},
	}

	rapid := Rapid{
		ctx:      context.Background(),
		teamName: teamName,
		log:      logger,
	}

	rapid.writer = &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			DialTimeout: 10 * time.Second,
			IdleTimeout: 10 * time.Second,
			TLS: &tls.Config{
				Certificates: []tls.Certificate{keypair},
				RootCAs:      caCertPool,
			},
		},
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:   []string{config.Brokers},
		Topic:     config.Topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
		Dialer:    dialer,
	}
	if config.GroupID != "" {
		readerConfig.GroupID = config.GroupID
	}

	rapid.reader = kafka.NewReader(readerConfig)

	return &rapid, nil
}

// Run starts the Rapid loop
func (r *Rapid) Run(handlers RapidHandlers) error {
	for {
		kafkaMessage, err := r.reader.FetchMessage(r.ctx)
		if err != nil {
			r.log.Error(fmt.Sprintf("failed to read message: %s", err))
			continue
		}

		var message LeesahMessage
		if err := json.Unmarshal(kafkaMessage.Value, &message); err != nil {
			r.log.Error(fmt.Sprintf("failed to unmarshal message: %s", err))
			continue
		}

		if message.TeamName == r.teamName {
			if err := r.handleTeamMessage(message, handlers); err != nil {
				r.log.Error(fmt.Sprintf("failed to handle message: %s", err))
				continue
			}
		}

		if err := r.reader.CommitMessages(r.ctx, kafkaMessage); err != nil {
			r.log.Error(fmt.Sprintf("failed to commit message: %s", err))
		}
	}
}

func (r *Rapid) handleTeamMessage(message LeesahMessage, handlers RapidHandlers) error {
	switch message.Type {
	case message_type_question:
		if handlers.HandleQuestion != nil {
			answer, err := handlers.HandleQuestion(message.ToQuestion())
			if err != nil {
				return fmt.Errorf("failed to handle question: %s", err)
			}

			if err := r.writeMessage(answer); err != nil {
				return fmt.Errorf("failed to write message: %s", err)
			}
		}
	case message_type_assessment:
		if handlers.HandleAssessment != nil {
			if err := handlers.HandleAssessment(message.ToAssessment()); err != nil {
				return fmt.Errorf("failed to handle assessment: %s", err)
			}
		}
	}

	return nil
}

// WriteMessage writes an Answer object to the Kafka topic. This function will automatically generate a MessageId and Created timestamp. It will also commit the offset of the message.
func (r *Rapid) writeMessage(answer Answer) error {
	message := LeesahMessage{
		Answer:     answer.Answer,
		AnswerID:   uuid.New(),
		MessageID:  uuid.New(),
		QuestionID: answer.QuestionID,
		Category:   answer.Category,
		Created:    time.Now().Format(time.RFC3339),
		TeamName:   r.teamName,
		Type:       message_type_answer,
	}

	output, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	return r.writer.WriteMessages(r.ctx, kafka.Message{Value: output})
}

// Close closes the Kafka writer and reader
func (r *Rapid) Close() {
	r.writer.Close()
	r.reader.Close()
}
