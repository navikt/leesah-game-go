# Go LEESAH

> Leesah-game is an event-driven application development game that challenges players to build an event-driven application.
> The application handles different types of tasks that it receives as events on a Kafka-based event stream.
> The tasks vary from very simple to more complex.

This is the Python library to play Leesah!

## Getting started

There are two versions of the the Leesah game!
On is for local play, directly in the terminal.
While the other is running on the Nais platform, and you learn how to be a developer in Nav and use Nais.
This library is used by both versions, but the following documentation is **just** for local play.


### Setting up local environment

You need Go to use this library, check out [go.dev/learn](https://go.dev/learn/) to get started with Go.

Start by creating the directory `leesah-game`.

```shell
mkdir leesah-game
cd leesah-game
go mod init leesah-game
touch main.go
```

### Install the library

There is only one dependency you need to play, and that is the [go-leesah](https://pkg.go.dev/github.com/navikt/go-leesah) library.

```shell
go get github.com/navikt/go-leesah
```

### Fetch Kafka certificates

You need some certificates to connect to the Kafka cluster, which is available at [leesah.io/certs](https://leesah.io/certs).
The username is always `leesah-game`, and the password will be distributed.

You can also use the one-liner below:

```bash
curl -u leesah-game:<see presentation> -o leesah-certs.zip https://leesah.io/certs && unzip leesah-certs.zip
```

Using the command above you will end up with `leesah-certs.yaml` in the `leesah-game` directory you made earlier.

### Example code

To make it easy to start we have made a working example that answer the first question, `team-registration`, with a dummy name and color.
All you need to do is update `TEAM_NAME` and `HEX_CODE`, and your ready to compete!

Copy the content below into `main.go`.
```go
package main

import (
	"log/slog"

	"github.com/navikt/go-leesah"
)

// 1. Ensure credential files are in the certs directory
// 2. Set `teamName` to your preferred team name
// 3. Set `teamColor` to your preferred team color

const (
    teamName  = "CHANGE ME"
    teamColor = "CHANGE ME"
)

func main() {
    log := slog.Default()
	ignoredCategories := []string{
		// "team-registration",
	}

	rapid, err := leesah.NewLocalRapid(teamName, log, ignoredCategories)
	if err != nil {
		log.Error("failed to create rapid", "error", err)
		return
	}
	defer rapid.Close()

	for {
		question, err := rapid.GetQuestion()
		if err != nil {
			slog.Error("can't get new question", "error", err)
		}

		var answer string
		switch question.Category {
		case "team-registration":
			answer = handleTeamRegistration(question)
		}

		if answer != "" {
			if err := rapid.Answer(answer); err != nil {
				log.Error("can't post answer", "error", err)
			}
		}
	}
}

func handleTeamRegistration(question leesah.Question) string {
	return teamColor
}
```

### Kjør lokalt

Så kan du kjøre det:

```shell
go run .
```
