# Go LEESAH

Go-biblitek for å spille LEESAH!

## Kom i gang

```go
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/navikt/go-leesah/pkg/leesah"
)

type Participant struct {
	log *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	os.Setenv("KAFKA_GROUP_ID", your-group-id)
	os.Setenv("KAFKA_TOPIC", "leesah-quiz.leesah-rapid-v2")

	rapid, err := leesah.NewRapid("my-team", logger)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create rapid: %s", err))
		return
	}
	defer rapid.Close()

	p := Participant{
		log: logger,
	}

	rapidHandlers := leesah.RapidHandlers{
		HandleQuestion: p.HandleQuestion,
	}

	if err := rapid.Run(rapidHandlers); err != nil {
		logger.Error(fmt.Sprintf("failed to run rapid: %s", err))
	}
}

func (p Participant) HandleQuestion(question leesah.Question) (leesah.Answer, error) {
	p.log.Info(fmt.Sprintf("%+v", question))

	# your code goes here

	return leesah.Answer{}, fmt.Errorf("no answer for this question")
}
```
### Lokal kjøring

```bash
nais aiven create --pool nav-dev kafka $USER leesah-quiz
# hent credentials med cmd fra output
```
