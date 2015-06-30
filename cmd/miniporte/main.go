package main

import (
	"os"
	"strings"

	bot "github.com/oz/miniporte"
)

func main() {
	bot.New(
		getEnvOr("IRC_SERVER", "irc.freenode.net:7000"),
		getEnvOr("IRC_NICK", "miniporte"),
		getEnvOr("IRC_NAME", "Mini-Porte"),
		getEnvOr("IRC_IDENT", "MiniPorteIRCBot"),
		strings.Split(getEnvOr("IRC_CHANS", "#af83-bots"), ","),
	).Run()
}

// Get the environment variable "name", or a default value.
func getEnvOr(name, defaultValue string) (out string) {
	out = os.Getenv(name)
	if out == "" {
		out = defaultValue
	}
	return
}
