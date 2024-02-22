# Go LEESAH

> Leesah-game er et hendelsedrevet applikasjonsutviklingspill som utfordrer spillerne til å bygge en hendelsedrevet applikasjon. 
> Applikasjonen håndterer forskjellige typer oppgaver som den mottar som hendelser på en Kafka-basert hendelsestrøm.
> Oppgavene varierer fra veldig enkle til mer komplekse.

Go-bibliotek for å spille LEESAH!

## Kom i gang

Det finnes to versjoner av Leesah-game!
En hvor man lager en applikasjon som kjører på Nais, og en hvor man spiller lokalt direkte fra terminalen sin.
Dette biblioteket kan brukes i begge versjoner, men denne dokumentasjonen dekker kun lokal spilling.
Vi har et eget template-repo som ligger under [navikt/leesah-game-template-go](https://github.com/navikt/leesah-game-template-go) for å spille Nais-versjonen.

### Hent credentials

Sertifikater for å koble seg på Kafka ligger tilgjengelig på [leesah-game-cert.ekstern.dev.nav.no/certs](https://leesah-game-cert.ekstern.dev.nav.no/certs), brukernavn og passord skal du få utdelt.
Du kan også bruke kommandoen nedenfor:

```bash
wget --user <username> --password <password> -O leesah-creds.zip https://leesah-game-cert.ekstern.dev.nav.no/certs && unzip leesah-creds.zip 
```

### Eksempelkode

Nedenfor er det et fungerende eksempel som svarer på lagregistreringsspørsmålet med et navn du velger, og en farge du velger:

```go
package main

import (
	"fmt"
	"log/slog"

	"github.com/navikt/go-leesah"
)

// 1. Ensure credential files are in the certs directory
// 2. Set `teamName` to your preferred team name
// 3. Set `teamColor` to your preferred team color

const (
    teamName  = "my-go-team"
    teamColor = "00ADD8"
)

func main() {
	rapid, err := leesah.NewLocalRapid(teamName, slog.Default())
	if err != nil {
		slog.Error("failed to create rapid", "error", err)
		return
	}
	defer rapid.Close()

	if err := rapid.Run(Answer); err != nil {
		slog.Error("failed to run rapid", "error", err)
	}
}

func Answer(question leesah.Question, log *slog.Logger) (string, bool) {
	slog.Info(fmt.Sprintf("%+v", question))

	switch question.Category {
	case "team-registration":
		return teamColor, true
	}

	return "", false
}
```

### Kjør lokalt

Først må du sette opp avhengigheten:

```shell
go get github.com/navikt/go-leesah
```

Så kan du kjøre det:

```shell
go run .
```
