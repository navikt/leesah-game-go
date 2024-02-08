# Go LEESAH

Go-bibliotek for å spille LEESAH!

## Kom i gang

Vi har et eget template-repo som ligger under [navikt/leesah-game-template-go](https://github.com/navikt/leesah-game-template-go).
Ellers kan du ta utgangspunkt i koden nedenfor.

```go
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/navikt/go-leesah"
)

var config leesah.RapidConfig

func init() {
	os.Setenv("KAFKA_GROUP_ID", uuid.New().String())
	os.Setenv("KAFKA_TOPIC", "leesah-quiz.leesah-rapid-v2")

	flag.StringVar(&config.Brokers, "brokers", os.Getenv("KAFKA_BROKERS"), "Kafka broker")
	flag.StringVar(&config.Topic, "topic", os.Getenv("KAFKA_TOPIC"), "Kafka topic")
	flag.StringVar(&config.GroupID, "group-id", os.Getenv("KAFKA_GROUP_ID"), "Kafka group ID")
	flag.StringVar(&config.KafkaCertPath, "kafka-cert-path", os.Getenv("KAFKA_CERTIFICATE_PATH"), "Path to Kafka certificate")
	flag.StringVar(&config.KafkaPrivateKeyPath, "kafka-private-key-path", os.Getenv("KAFKA_PRIVATE_KEY_PATH"), "Path to Kafka private key")
	flag.StringVar(&config.KafkaCAPath, "kafka-ca-path", os.Getenv("KAFKA_CA_PATH"), "Path to Kafka CA certificate")
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config.Log = log

	rapid, err := leesah.NewRapid("my-go-team", config)
	if err != nil {
		log.Error(fmt.Sprintf("failed to create rapid: %s", err))
		return
	}
	defer rapid.Close()

	if err := rapid.Run(Answer); err != nil {
		log.Error(fmt.Sprintf("failed to run rapid: %s", err))
	}
}

func Answer(question leesah.Question, log *slog.Logger) (string, bool) {
	log.Info(fmt.Sprintf("%+v", question))

	var answer string
	switch question.Category {
	case "team-registration":
		return "", true
	}

	log.Info(fmt.Sprintf("{Category:%s Question:%s Answer:%s", question.Category, question.Question, answer))
	return "", false
}

```
### Lokal kjøring

```bash
nais aiven create --pool nav-dev kafka $USER leesah-quiz
# hent credentials med cmd fra output
```
