# Miniporte

This is a small IRC bot to save links shared on IRC channels onto various
online services.

Every link is logged using the provided tags. When no tag is provided, or if
the *#private* tag is used, the link is saved privately (for services that
support it).

Additionally, the IRC channel, and user's nick are used as tags too.

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
- `IRC_CHANS` comma-separated list of IRC channels, defaults to `#minibots`.

In order to use [Epistoli][epistoli]'s cute API to build newsletters
from the links shared on IRC, you will also need:

- `EPISTOLI_TOKEN` API token to post links,
- `EPISTOLI_LETTER` The newsletter name where links are posted.

[epistoli]: https://episto.li

# Usage

Just launch the `miniporte` program if the defaults are fine. Check the `-help`
flag for more options.

# Todo

Lots. :muscle:

# Background

This is maybe boring, TLDR: a simpler version of another proprietary IRC bot.

It all started at [af83][af83], where a [Cinch][cinch]-based bot would log
links from IRC to a Delicious account (for later processing). The bot did many
things, talked to Jenkins, Trello, etc.

Well it did a lot of other (possibly) useful things in those days of hardware
hacking, and software exploration, and fun. It was called [Mr. Porte][mr-porte]
(*Mr. Door* opened the electronic front-door lock), and it is now offline.

This project, however, is content with the logging of links (URIs really) onto
your favourite online service. Hence its name *miniporte*, tiny-door, as a
tribute to a long gone ancestor. :door:

[af83]: http://af83.com
[cinch]: https://github.com/cinchrb/cinch
[mr-porte]: https://episto.li/profile/les-pepitos/mr-porte

# License

MIT.
