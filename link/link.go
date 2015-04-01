package link

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	UserAgent        = "Opendoor"
	DeliciousPostUrl = "https://api.del.icio.us/v1/posts/add"
	HttpTimeout      = 5 // Wait at most 5 seconds for Delicious.
)

var (
	LinkRx = regexp.MustCompile(`(https?://[^ )]+)`)
	TagsRx = regexp.MustCompile(`\s(#[\w\pL-]+)`)
	CodeRx = regexp.MustCompile(`code="(.*?)"`)
)

func Find(txt string) string {
	return LinkRx.FindString(txt)
}

// Extract tags from a string, w/o the leading '#'
func Tags(txt string) []string {
	tags := []string{}
	for _, match := range TagsRx.FindAllStringSubmatch(txt, -1) {
		tags = append(tags, strings.TrimLeft(match[1], "#"))
	}
	return tags
}

func deliciousPostParams(u string, tags []string) io.Reader {
	form := url.Values{}
	form.Set("url", u)
	form.Set("description", u) // FIXME actually a title
	form.Set("tags", strings.Join(tags, ","))
	if IncludesPrivate(tags) {
		form.Set("shared", "no")
	}
	return strings.NewReader(form.Encode())
}

// Custom Delicious API client, without keep-alive.
func deliciousClient() *http.Client {
	return &http.Client{
		Timeout:   HttpTimeout * time.Second,
		Transport: &http.Transport{DisableKeepAlives: true},
	}
}

func Save(u string, tags []string) (err error) {
	params := deliciousPostParams(u, tags)
	req, err := http.NewRequest("POST", DeliciousPostUrl, params)
	if err != nil {
		return
	}
	token, err := oauthToken()
	if err != nil {
		return
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := deliciousClient()
	if resp, err := client.Do(req); err == nil {
		return parseResponse(resp)
	}
	return
}

func IncludesPrivate(tags []string) bool {
	for _, t := range tags {
		if t == "#private" || t == "private" {
			return true
		}
	}
	return false
}

// Delicious has a shitty API. Forget about HTTP status codes, embrace XML!
func parseResponse(resp *http.Response) error {
	defer resp.Body.Close()

	// We need to read code's value at <result code="value" />, to acknowledge
	// the API's response: HTTP codes are probably too hard.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Delicious hung up...")
	}
	m := CodeRx.FindSubmatch(body)
	if len(m) != 2 {
		return errors.New(string(body))
	}
	code := string(m[1])

	// "done" is good, we like "done".
	if code == "done" {
		return nil
	}
	log.Println("Delicious error: ", string(body))
	return errors.New(code)
}

func oauthToken() (string, error) {
	token := os.Getenv("DELICIOUS_OAUTH_TOKEN")
	if token == "" {
		return "", errors.New("Missing Delicious Token")
	}
	return "Bearer " + token, nil
}
