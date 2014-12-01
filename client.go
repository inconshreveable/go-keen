package keen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseUrl = "https://api.keen.io/3.0/projects/"
)

type KeenProperties struct {
	Timestamp string `json:"timestamp"`
}

// Timestamp formats a time.Time object in the ISO-8601 format keen expects
func Timestamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

type Client struct {
	WriteKey   string
	ProjectID  string
	HttpClient http.Client
}

func (c *Client) AddEvent(collection string, event interface{}) error {
	resp, err := c.request("POST", fmt.Sprintf("/events/%s", collection), event)
	if err != nil {
		return err
	}

	return c.respToError(resp)
}

func (c *Client) AddEvents(events map[string][]interface{}) error {
	resp, err := c.request("POST", "/events", events)
	if err != nil {
		return err
	}

	return c.respToError(resp)
}

func (c *Client) respToError(resp *http.Response) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("Non 200 reply from keen.io: %s", data)
}

func (c *Client) request(method, path string, payload interface{}) (*http.Response, error) {
	// serialize payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// construct url
	url := baseUrl + c.ProjectID + path

	// new request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// add auth
	req.Header.Add("Authorization", c.WriteKey)

	// set length/content-type
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.ContentLength = int64(len(body))
	}

	return c.HttpClient.Do(req)
}
