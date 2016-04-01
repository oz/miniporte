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
	// ServiceName is this package's pretty name.
	ServiceName = "epistoli"

	// APIURL is the base API endpoint for Epistoli.
	APIURL = "https://episto.li/api/v1"

	// HTTPTimeout specifies how many seconds we wait for the API to answer.
	HTTPTimeout = 3
)

// Epistoli base type: just enough to satisfy the link.Service
// interface.
type Epistoli struct {
	Letter string
	client *http.Client
}

type response struct {
	Ok  bool   `json:"ok"`
	Err string `json:"error"`
	// And other stuff we don't care about ATM...
}

// New initialize an Epistoli service.
func New() *Epistoli {
	return &Epistoli{
		Letter: os.Getenv("EPISTOLI_LETTER"),
		client: &http.Client{
			Timeout:   HTTPTimeout * time.Second,
			Transport: &http.Transport{DisableKeepAlives: true},
		},
	}
}

func (e *Epistoli) String() string {
	return ServiceName
}

// Save a Link to Epistoli
func (e *Epistoli) Save(l *link.Link) (err error) {
	params := e.postParams(l)
	postURL := APIURL + "/bookmarks"

	req, err := http.NewRequest("POST", postURL, params)
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
func (e *Epistoli) postParams(l *link.Link) io.Reader {
	form := url.Values{}
	form.Set("url", l.Url)
	form.Set("newsletter_id", e.Letter)
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
