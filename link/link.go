package link

import (
	"errors"
	"regexp"
	"strings"
)

const (
	UserAgent   = "Miniporte IRC bot 0.1"
	NoLinkError = "No link found"
)

var (
	LinkRx = regexp.MustCompile(`(https?://[^ )]+)`)
	TagsRx = regexp.MustCompile(`\s(#[\w\pL-]+)`)
)

type Service interface {
	String() string
	Save(*Link) error
}

type Link struct {
	Url     string
	Tags    []string
	Pub     bool
	Service Service
}

func New(s Service) *Link {
	return &Link{
		Service: s,
	}
}

// Extract a Link struct from a string
func (l *Link) MustExtract(txt string) error {
	url := LinkRx.FindString(txt)
	if url == "" {
		return errors.New(NoLinkError)
	}
	tags := extractTags(txt)
	l.Tags = tags
	l.Url = url
	l.Pub = isPublic(tags)

	return nil
}

func (l *Link) Save() (err error) {
	return l.Service.Save(l)
}

// Extract tags from a string, w/o the leading '#'
func extractTags(txt string) []string {
	tags := []string{}
	for _, match := range TagsRx.FindAllStringSubmatch(txt, -1) {
		tags = append(tags, strings.TrimLeft(match[1], "#"))
	}
	return tags
}

func isPublic(tags []string) bool {
	if len(tags) == 0 {
		return false
	}
	for _, t := range tags {
		if t == "#private" || t == "private" {
			return false
		}
	}
	return true
}
