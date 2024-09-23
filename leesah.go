package leesah

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

// rapid is a struct that represents a connection to the Leesah Kafka topic
// Use the getQuestion method to get the next question from the topic, and
// the answer method to post your answer to the topic.
type rapid struct {
	ignoredCategories []string
	writer            *kafka.Writer
	reader            *kafka.Reader
	ctx               context.Context
	teamName          string
	lastMessage       *Message
	log               *slog.Logger
	kafkaDir          string
}

// RapidConfig is a struct that represents the configuration for a Rapid instance
// This is used when creating a new Rapid instance on the NAIS platform
type RapidConfig struct {
	Broker            string
	CAPath            string
	CertPath          string
	GroupID           string
	IgnoredCategories []string
	KafkaDir          string
	Log               *slog.Logger
	PrivateKeyPath    string
	Topic             string
}

// NewLocalRapid creates a new Rapid instance with a local configuration.
// The local configuration is read from "leesah-certs.yaml".
// It is used when playing the local edition of Leesah.
// You can override the path to the local certification by setting the
// environment variable QUIZ_CERTS.
// You can also override the topic by setting the environment variable
// QUIZ_TOPIC, or else the first topic in the file will be used.
func NewLocalRapid(teamName string, log *slog.Logger, ignoredCategories []string) (*rapid, error) {
	rapidConfig, err := loadLocalConfig(log, ignoredCategories)
	if err != nil {
		return nil, fmt.Errorf("failed to load local config: %s", err)
	}

	rapid, err := NewRapid(teamName, rapidConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create rapid: %s", err)
	}

	return rapid, nil
}

func loadLocalConfig(log *slog.Logger, ignoredCategories []string) (RapidConfig, error) {
	log.Info("⚙️ Loading local config")

	certPath := "leesah-certs.yaml"
	if os.Getenv("QUIZ_CERTS") != "" {
		certPath = os.Getenv("QUIZ_CERTS")
	}

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		certPath = "certs/leesah-certs.yaml"
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			return RapidConfig{}, fmt.Errorf("can't find certs in 'leesah-certs.yaml', '$QUIZ_CERTS', or 'certs/leesah-certs.yaml'")
		}
	}

	localFile, err := os.ReadFile(certPath)
	if err != nil {
		return RapidConfig{}, fmt.Errorf("failed to read local file: %s", err)
	}

	type creds struct {
		Broker string   `yaml:"broker"`
		Topics []string `yaml:"topics"`
		CA     string   `yaml:"ca"`
		User   struct {
			AccessKey         string `yaml:"access_key"`
			AccessCertificate string `yaml:"access_cert"`
			Username          string `yaml:"username"`
		} `yaml:"user"`
	}

	var c creds
	if err := yaml.Unmarshal(localFile, &c); err != nil {
		return RapidConfig{}, fmt.Errorf("failed to unmarshal local file: %s", err)
	}

	dir, err := os.MkdirTemp("", "kafka")
	if err != nil {
		return RapidConfig{}, err
	}

	caFile, err := writeToTempDir(dir, "ca.pem", c.CA)
	if err != nil {
		return RapidConfig{}, err
	}

	certFile, err := writeToTempDir(dir, "cert.crt", c.User.AccessCertificate)
	if err != nil {
		return RapidConfig{}, err
	}

	privateKeyFile, err := writeToTempDir(dir, "private-key.pem", c.User.AccessKey)
	if err != nil {
		return RapidConfig{}, err
	}

	topic := c.Topics[0]
	if os.Getenv("QUIZ_TOPIC") != "" {
		topic = os.Getenv("QUIZ_TOPIC")
	}

	if !strings.Contains(c.Broker, ":") {
		c.Broker += ":26484"
	}

	return RapidConfig{
		Log:               log,
		Broker:            c.Broker,
		Topic:             topic,
		GroupID:           uuid.New().String(),
		CertPath:          certFile,
		PrivateKeyPath:    privateKeyFile,
		CAPath:            caFile,
		KafkaDir:          dir,
		IgnoredCategories: ignoredCategories,
	}, nil
}

func writeToTempDir(dir, fileName, data string) (string, error) {
	filePath := filepath.Join(dir, fileName)
	if err := os.WriteFile(filePath, []byte(data), 0o666); err != nil {
		return "", err
	}

	return filePath, nil
}

// NewRapid creates a new Rapid instance with the given configuration.
// It is used when playing the Nais-edition of Leesah.
func NewRapid(teamName string, config RapidConfig) (*rapid, error) {
	config.Log.Info("🔨 Creating new rapid")
	keypair, err := tls.LoadX509KeyPair(config.CertPath, config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load Access Key and/or Access Certificate: %s", err)
	}

	caCert, err := os.ReadFile(config.CAPath)
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

	rapid := rapid{
		ctx:               context.Background(),
		teamName:          teamName,
		ignoredCategories: config.IgnoredCategories,
		log:               config.Log,
		kafkaDir:          config.KafkaDir,
	}

	rapid.writer = &kafka.Writer{
		Addr:     kafka.TCP(config.Broker),
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
		Brokers:   []string{config.Broker},
		Topic:     config.Topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
		Dialer:    dialer,
		GroupID:   config.GroupID,
	}

	rapid.reader = kafka.NewReader(readerConfig)

	rapid.log.Info("🚀 Starting QuizRapid...")
	rapid.log.Info("🔍 Looking for first question")

	return &rapid, nil
}

// GetQuestion will wait for the next question from the Kafka topic
// and return it as a Question struct.
func (r *rapid) GetQuestion() (Question, error) {
	for {
		kafkaMessage, err := r.reader.FetchMessage(r.ctx)
		if err != nil {
			r.log.Error(fmt.Sprintf("failed to read message: %s", err))
			continue
		}

		var mm minimalMessage
		if err := json.Unmarshal(kafkaMessage.Value, &mm); err != nil {
			r.log.Debug(string(kafkaMessage.Value))
			return Question{}, fmt.Errorf("failed to unmarshal minimal message: %s", err)
		}

		// defer r.reader.CommitMessages(r.ctx, kafkaMessage)

		if mm.Type == MessageTypeQuestion {
			var message Message
			if err := json.Unmarshal(kafkaMessage.Value, &message); err != nil {
				r.log.Debug(string(kafkaMessage.Value))
				return Question{}, fmt.Errorf("failed to unmarshal message: %s", err)
			}

			r.lastMessage = &message
			question := message.ToQuestion()
			if !slices.Contains(r.ignoredCategories, question.Category) {
				r.log.Info(fmt.Sprintf("📥 Received question: kategorinavn='%s' spørsmål='%s' svarformat='%s' id='%s' dokumentasjon='%s'", question.Category, question.Question, question.AnswerFormat, question.ID, question.Documentation))
			}

			return question, nil
		}
	}
}

// Answer creates a Kafka message and posts it to the Kafka topic
func (r *rapid) Answer(answer string) error {
	kafkaMessage := Message{
		Answer:     answer,
		Category:   r.lastMessage.Category,
		Created:    time.Now().Format(LeesahTimeformat),
		AnswerID:   uuid.New().String(),
		QuestionID: r.lastMessage.QuestionID,
		TeamName:   r.teamName,
		Type:       MessageTypeAnswer,
	}

	output, err := json.Marshal(kafkaMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	if err := r.writer.WriteMessages(r.ctx, kafka.Message{Value: output}); err != nil {
		return fmt.Errorf("failed to write message: %s", err)
	}

	if !slices.Contains(r.ignoredCategories, r.lastMessage.Category) {
		r.log.Info(fmt.Sprintf("📤 Published answer: kategorinavn='%s' svar='%s' lagnavn='%s'", r.lastMessage.Category, answer, r.teamName))
	}

	r.lastMessage = nil

	return nil
}

// Close closes the Kafka writer and reader
func (r *rapid) Close() {
	r.writer.Close()
	r.reader.Close()
	if r.kafkaDir != "" {
		os.RemoveAll(r.kafkaDir)
	}
}
