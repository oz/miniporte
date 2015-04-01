# Miniporte

This is a small IRC bot we use at [af83](http://af83.com) to log links to a
Delicious account (before processing them).  It does not do anything appart
from that.

Every link is logged using the provided tags. When no tag is provided, or if
the *#private* tag is used, the link is saved privately.

# Installation

```
$ go get github.com/oz/miniporte
$ go install github.com/oz/miniporte/cmd/miniporte
```

# Configuration

The bot is configured through the following environment variables.

- `IRC_SERVER` IRC server, defaults to `irc.freenode.net:7000` (SSL is
  *always* on).
- `IRC_NICK` IRC nick, defaults to `miniporte`.
- `IRC_NAME` IRC name, defaults to `Mini-Porte`.
- `IRC_IDENT` IRC *ident* name, defaults to `MiniPorteIRCBot`.
- `IRC_CHANS` comma-separated list of IRC channels, defaults to
  `#af83-bots`.
- `DELICIOUS_OAUTH_TOKEN` we use Delicious' API (for the worse), get an
  OAuth token to post links, set it here.

# Usage

Just launch the `miniporte` binary, it has no command-line flags.

# Todo

Lots.

# License

MIT.
