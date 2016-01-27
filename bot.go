package miniporte

import (
	"log"

	"github.com/cenkalti/backoff"
	irc "github.com/fluffle/goirc/client"
	"github.com/oz/miniporte/epistoli"
	link "github.com/oz/miniporte/link"
)

type Bot struct {
	Chans  []string
	Config *irc.Config
	Client *irc.Conn
	Ctl    (chan string)
}

// Create a new Bot
func New(server, nick, name, ident string, chans []string) *Bot {
	cfg := irc.NewConfig(nick)
	cfg.SSL = true
	cfg.Me.Name = name
	cfg.Me.Ident = ident
	cfg.Server = server
	cfg.NewNick = func(n string) string { return n + "_" }

	return &Bot{
		Chans:  chans,
		Config: cfg,
		Client: irc.Client(cfg),
		Ctl:    make(chan string),
	}
}

func (b *Bot) joinChannels() {
	log.Println("Joining channels", b.Chans)
	for _, c := range b.Chans {
		b.Client.Join(c)
	}
}

func (b *Bot) Run() {
	b.initializeHandlers()
	go b.Connect()
	b.commandLoop()
	log.Println("Bot quitting...")
}

func (b *Bot) onMessage(msg *irc.Line) {
	// Ignore non-public messages
	if !msg.Public() {
		return
	}

	l := link.New(epistoli.New())
	if err := l.MustExtract(msg.Text()); err != nil {
		// No links.
		return
	}
	l.Tags = append(l.Tags, msg.Nick, msg.Target())
	go func() {
		if err := l.Save(); err != nil {
			if l.Pub {
				b.Client.Privmsg(msg.Target(), "Oops! "+err.Error())
			}
			return
		}
		if l.Pub {
			b.Client.Privmsg(msg.Target(), "Saved!")
		}
	}()
}

func (b *Bot) initializeHandlers() {
	// Connected
	b.Client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		log.Println("Connected!")
		b.joinChannels()
	})

	// Disconnected
	b.Client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		log.Println("Disconnected")
		b.Ctl <- "disconnected"
	})

	// PRIVMSG
	b.Client.HandleFunc("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		b.onMessage(line)
	})
}

func (b *Bot) Connect() error {
	log.Println("Connecting to", b.Config.Server)
	return b.Client.Connect()
}

func (b *Bot) connect() {
	if err := backoff.Retry(b.Connect, backoff.NewExponentialBackOff()); err != nil {
		log.Printf("Connection error: %s\n", err)
		b.Ctl <- "connection-error"
	}
}

// Connection loop
func (b *Bot) commandLoop() {
	for {
		for cmd := range b.Ctl {
			switch cmd {
			case "disconnected", "connection-error":
				log.Println("Disconnected:", cmd)
				go b.connect()
			default:
				log.Println("Ignoring command", cmd)
			}
		}
	}
}
