# Go LEESAH

> Leesah-game er et hendelsedrevet applikasjonsutviklingspill som utfordrer spillerne til å bygge en hendelsedrevet applikasjon. 
> Applikasjonen håndterer forskjellige typer oppgaver som den mottar som hendelser på en Kafka-basert hendelsestrøm.
> Oppgavene varierer fra veldig enkle til mer komplekse.

Go-bibliotek for å spille LEESAH!

## Kom i gang

Det finnes to versjoner av Leesah Game!
En hvor man lager en applikasjon som kjører på Nais, og en hvor man spiller lokalt direkte fra terminalen sin.
Dette biblioteket kan brukes i begge versjoner, men denne dokumentasjonen dekker **kun** lokal spilling.

Vi har et eget template-repo som ligger under [navikt/leesah-game-template-go](https://github.com/navikt/leesah-game-template-go) for å spille Nais-versjonen.

### Sett opp lokalt miljø

Du trenger Go for å bruke biblioteket, sjekk ut [go.dev/learn](https://go.dev/learn/) for å komme i gang.

Start med å opprette en katalog `leesah-game`.

```shell
mkdir leesah-game
cd leesah-game
go mod init leesah-game
touch main.go
```

### Installer biblioteket

Det er kun en ekstern avhengighet du trenger, og det er biblioteket [go-leesah](https://pkg.go.dev/github.com/navikt/go-leesah).

```shell
go get github.com/navikt/go-leesah
```

### Hent Kafkasertifakt

Sertifikater for å koble seg på Kafka ligger tilgjengelig på [leesah.io/certs](https://leesah.io/certs), brukernavn er `leesah-game`, passord bli utdelt.

Du kan også bruke kommandoen nedenfor:

```bash
curl -u leesah-game:<passord> -o leesah-certs.zip https://leesah.io/certs && unzip leesah-certs.zip
```

### Eksempelkode

Nedenfor er det et fungerende eksempel som du kan lime inn i `main.go` som svarer på spørsmålet om lagregistrerings med et navn, og en farge (hex-kode):

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
		// "lagregistrering",
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
		case "lagregistrering":
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
