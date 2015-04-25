package miniporte

import (
	"log"
	"os"
	"strings"

	irc "github.com/fluffle/goirc/client"
	link "github.com/oz/miniporte/link"
)

type Bot struct {
	Chans  []string
	Config *irc.Config
	Client *irc.Conn
	Ctl    (chan string)
}

func New() *Bot {
	cfg := irc.NewConfig(getEnvOr("IRC_NICK", "miniporte"))
	cfg.SSL = true
	cfg.Me.Name = getEnvOr("IRC_NAME", "Mini-Porte")
	cfg.Me.Ident = getEnvOr("IRC_IDENT", "MiniPorteIRCBot")
	cfg.Server = getEnvOr("IRC_SERVER", "irc.freenode.net:7000")
	cfg.NewNick = func(n string) string { return n + "_" }

	return &Bot{
		Chans:  strings.Split(getEnvOr("IRC_CHANS", "#af83-bots"), ","),
		Config: cfg,
		Client: irc.Client(cfg),
		Ctl:    make(chan string),
	}
}

func (b *Bot) OnMessage(msg *irc.Line) {
	// Ignore non-public messages
	if !msg.Public() {
		return
	}

	if url := link.Find(msg.Text()); url != "" {
		tags := link.Tags(msg.Text())
		if len(tags) == 0 {
			tags = []string{"private"}
		}
		tags = append(tags, msg.Nick, msg.Target())
		go func() {
			if err := link.Save(url, tags); err != nil {
				if !link.IncludesPrivate(tags) {
					b.Client.Privmsg(msg.Target(), "Oops! "+err.Error())
				}
				return
			}
			if !link.IncludesPrivate(tags) {
				b.Client.Privmsg(msg.Target(), "Saved!")
			}
		}()
	}
}

func (b *Bot) JoinChannels() {
	log.Println("Joining channels", b.Chans)
	for _, c := range b.Chans {
		b.Client.Join(c)
	}
}

func (b *Bot) Run() {
	// Connected
	b.Client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		log.Println("Connected!")
		b.JoinChannels()
	})

	// Disconnected
	b.Client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		log.Println("Disconnected")
		b.Ctl <- "disconnected"
	})

	// PRIVMSG
	b.Client.HandleFunc("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		b.OnMessage(line)
	})

	// Connection loop
	for {
		log.Println("Connecting to IRC...")
		if err := b.Client.Connect(); err != nil {
			log.Printf("Connection error: %s\n", err)
		}

		for cmd := range b.Ctl {
			switch cmd {
			case "quit":
				b.Client.Quit("Bye...")
				log.Println("Quitting...")
				return
			case "disconnected":
				log.Println("Trying to reconnect after", cmd)
				break
			default:
				log.Println("Ignoring command", cmd)
			}
		}
	}
}

// Retrieve the environment variable "name", or a default value.
func getEnvOr(name, defaultValue string) (out string) {
	out = os.Getenv(name)
	if out == "" {
		out = defaultValue
	}
	return
}
