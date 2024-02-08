package leesah

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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

type RapidConfig struct {
	Brokers             string
	Topic               string
	GroupID             string
	KafkaCertPath       string
	KafkaPrivateKeyPath string
	KafkaCAPath         string
	Log                 *slog.Logger
}

func NewRapid(teamName string, config RapidConfig) (*Rapid, error) {
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
		log:      config.Log,
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

	if config.GroupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:   []string{config.Brokers},
		Topic:     config.Topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
		Dialer:    dialer,
		GroupID:   config.GroupID,
	}

	rapid.reader = kafka.NewReader(readerConfig)

	return &rapid, nil
}

func (r *Rapid) Run(answerQuestion func(Question, *slog.Logger) (string, bool)) error {
	for {
		kafkaMessage, err := r.reader.FetchMessage(r.ctx)
		if err != nil {
			r.log.Error(fmt.Sprintf("failed to read message: %s", err))
			continue
		}

		var message Message
		if err := json.Unmarshal(kafkaMessage.Value, &message); err != nil {
			return fmt.Errorf("failed to unmarshal message: %s", err)
		}

		answer, ok := answerQuestion(message.ToQuestion(), r.log)
		if !ok {
			continue
		}

		if err := r.postAnswer(message, answer); err != nil {
			return fmt.Errorf("failed to post answer: %s", err)
		}

		if err := r.reader.CommitMessages(r.ctx, kafkaMessage); err != nil {
			return fmt.Errorf("failed to commit message: %s", err)
		}
	}
}

func (r *Rapid) postAnswer(message Message, answer string) error {
	kafkaMessage := Message{
		MessageID:  uuid.New(),
		QuestionID: message.MessageID,
		Category:   message.Category,
		Created:    time.Now().Format(LeesahTimeformat),
		TeamName:   r.teamName,
		Type:       MessageTypeAnswer,
		Answer:     answer,
	}

	return r.writeMessage(kafkaMessage)
}

func (r *Rapid) writeMessage(message Message) error {
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
