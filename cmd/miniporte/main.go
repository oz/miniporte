package main

import (
	"flag"
	"os"
	"strings"

	bot "github.com/oz/miniporte"
)

func main() {
	var (
		server = flag.String("server", getEnvOr("IRC_SERVER", "chat.freenode.net:7000"), "IRC Server")
		nick   = flag.String("nick", getEnvOr("IRC_NICK", "miniporte"), "Bot nick")
		name   = flag.String("name", getEnvOr("IRC_NAME", "Mini-Porte"), "Bot's name")
		ident  = flag.String("ident", getEnvOr("IRC_IDENT", "MiniPorteIRCBot"), "Bot's ident")
		chans  = flag.String("chans", getEnvOr("IRC_CHANS", "#minibots"), "Bot's chans at boot, comma separated")
	)
	flag.Parse()
	bot.New(*server, *nick, *name, *ident, strings.Split(*chans, ",")).Run()
}

// Get the environment variable "name", or a default value.
func getEnvOr(name, defaultValue string) (out string) {
	out = os.Getenv(name)
	if out == "" {
		out = defaultValue
	}
	return
}
