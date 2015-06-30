package miniporte

import (
	"log"

	irc "github.com/fluffle/goirc/client"
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
	b.commandLoop()
	log.Println("Bot quitting...")
}

func (b *Bot) onMessage(msg *irc.Line) {
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

// Connection loop
func (b *Bot) commandLoop() {
	for {
		log.Println("Connecting to IRC...")
		if err := b.Client.Connect(); err != nil {
			log.Printf("Connection error: %s\n", err)
		}

		for cmd := range b.Ctl {
			switch cmd {
			case "quit":
				b.Client.Quit("Bye...")
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
