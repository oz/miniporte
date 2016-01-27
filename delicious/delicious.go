package delicious

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

	"github.com/oz/miniporte/link"
)

const (
	DeliciousPostUrl = "https://api.del.icio.us/v1/posts/add"
	HttpTimeout      = 5 // Wait at most 5 seconds for Delicious.
	ServiceName      = "delicious"
)

var (
	CodeRx = regexp.MustCompile(`code="(.*?)"`)
)

type Delicious struct {
	client *http.Client
}

// Custom Delicious API client, without keep-alive.
func New() *Delicious {
	return &Delicious{
		client: &http.Client{
			Timeout:   HttpTimeout * time.Second,
			Transport: &http.Transport{DisableKeepAlives: true},
		},
	}
}

func (d *Delicious) String() string {
	return ServiceName
}

func (d *Delicious) Save(l *link.Link) (err error) {
	params := postParams(l)
	req, err := http.NewRequest("POST", DeliciousPostUrl, params)
	if err != nil {
		return
	}
	token, err := oauthToken()
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)
	req.Header.Add("User-Agent", link.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err := d.client.Do(req); err == nil {
		return parseResponse(resp)
	}
	return
}

func postParams(l *link.Link) io.Reader {
	form := url.Values{}
	form.Set("url", l.Url)
	form.Set("description", l.Url) // FIXME actually a title
	form.Set("tags", strings.Join(l.Tags, ","))
	if !l.Pub {
		form.Set("shared", "no")
	}
	return strings.NewReader(form.Encode())
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
	if code == "done" || code == "error adding link" {
		return nil
	}
	log.Println("Delicious error: ", string(body))
	return errors.New(code)
}

func oauthToken() (string, error) {
	token := os.Getenv("DELICIOUS_OAUTH_TOKEN")
	if token == "" {
		return "", errors.New("Missing Delicious OAuth Token")
	}
	return "Bearer " + token, nil
}
