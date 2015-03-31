package link

import (
	"errors"
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
	HttpTimeout      = 30 // 30 seconds timeout, yup.
)

var (
	LinkRx = regexp.MustCompile(`(https?://[^ )]+)`)
	TagsRx = regexp.MustCompile(`#[\w-]+`)
	CodeRx = regexp.MustCompile(`code="(.*?)"`)
)

func Find(txt string) string {
	return LinkRx.FindString(txt)
}

// Extract tags from a string, w/o the leading '#'
func Tags(txt string) []string {
	tags := []string{}
	for _, tag := range TagsRx.FindAllString(txt, -1) {
		tags = append(tags, strings.TrimLeft(tag, "#"))
	}
	return tags
}

func Save(u string, tags []string) (err error) {
	client := &http.Client{Timeout: HttpTimeout * time.Second}
	params := url.Values{}
	params.Set("url", u)
	params.Set("description", u) // FIXME actually a title
	params.Set("tags", strings.Join(tags, ","))
	if IncludesPrivate(tags) {
		params.Set("shared", "no")
	}

	req, err := http.NewRequest("POST", DeliciousPostUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}
	token, err := oauthToken()
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	return parseResponse(resp)
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
