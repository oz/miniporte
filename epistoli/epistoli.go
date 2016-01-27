package epistoli

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/oz/miniporte/link"
)

const (
	ServiceName = "epistoli"
	ApiUrl      = "https://episto.li/api/v1"
	HttpTimeout = 3 // no patience, 3 secs is plenty!
)

// The basic Epistoli type: just enough to satisfy the link.Service
// interface.
type Epistoli struct {
	client *http.Client
}

type response struct {
	Ok  bool   `json:"ok"`
	Err string `json:"error"`
	// And other stuff we don't care about ATM...
}

func New() *Epistoli {
	return &Epistoli{
		client: &http.Client{
			Timeout:   HttpTimeout * time.Second,
			Transport: &http.Transport{DisableKeepAlives: true},
		},
	}
}

func (e *Epistoli) String() string {
	return ServiceName
}

// Save a Link to Epistoli
func (e *Epistoli) Save(l *link.Link) (err error) {
	params := postParams(l)
	postUrl := ApiUrl + "/bookmarks"

	req, err := http.NewRequest("POST", postUrl, params)
	if err != nil {
		return
	}
	token, err := authToken()
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Token", token)
	req.Header.Add("User-Agent", link.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err := e.client.Do(req); err == nil {
		return parseResponse(resp)
	}
	return
}

func parseResponse(resp *http.Response) error {
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	jsonRes := response{}
	err := d.Decode(&jsonRes)
	// Request error?
	if err != nil {
		return err
	}
	// Epistoli error?
	if !jsonRes.Ok {
		return errors.New(jsonRes.Err)
	}
	// All good.
	return nil
}

// Check the docs at: https://github.com/Epistoli/apidocs
func postParams(l *link.Link) io.Reader {
	form := url.Values{}
	form.Set("url", l.Url)
	form.Set("tags", strings.Join(l.Tags, ","))
	return strings.NewReader(form.Encode())
}

func authToken() (string, error) {
	token := os.Getenv("EPISTOLI_TOKEN")
	if token == "" {
		return "", errors.New("Missing Epistoli Token")
	}
	return token, nil
}
