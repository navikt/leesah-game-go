# Go LEESAH

Go-biblitek for å spille LEESAH!

## Kom i gang

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

type Participant struct {
	log *slog.Logger
}

var config = leesah.RapidConfig{}

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

	rapid, err := leesah.NewRapid("my-team", config)
	if err != nil {
		log.Error(fmt.Sprintf("failed to create rapid: %s", err))
		return
	}
	defer rapid.Close()

	p := Participant{
		log: log,
	}

	rapidHandlers := leesah.RapidHandlers{
		HandleQuestion: p.HandleQuestion,
	}

	if err := rapid.Run(rapidHandlers); err != nil {
		log.Error(fmt.Sprintf("failed to run rapid: %s", err))
	}
}

func (p Participant) HandleQuestion(question leesah.Question) (string, bool) {
	p.log.Info(fmt.Sprintf("%+v", question))

	var answer string
	switch question.Category {
	case "team-registration":
		answer = "00ADD8"
	default:
		p.log.Error(fmt.Sprintf("no answer for %s question", question.Category))
		return "", false
	}

	p.log.Info(fmt.Sprintf("{Category:%s Question:%s Answer:%s", question.Category, question.Question, answer))
	return answer, true
}
```
### Lokal kjøring

```bash
nais aiven create --pool nav-dev kafka $USER leesah-quiz
# hent credentials med cmd fra output
```
