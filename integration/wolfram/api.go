package wolfram

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL = "https://api.wolframalpha.com/v1/"

// Client stores AppID and also has API method
type Client struct {
	AppID string
}

var (
	// ErrInvalidInput thrown when user input is invalid
	ErrInvalidInput = errors.New("invalid input is given")
	// ErrInvalidAppID thrown when the given App ID is invalid
	ErrInvalidAppID = errors.New("invalid appid")
)

func (c *Client) buildURL(uri string, params map[string]string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	base.Path += uri
	baseParams := url.Values{}
	baseParams.Add("appid", c.AppID)
	for k, v := range params {
		baseParams.Add(k, v)
	}
	base.RawQuery = baseParams.Encode()
	return base.String(), nil
}

// Simple is an API for WolframAlpha's simple result image
func (c *Client) Simple(input string) ([]byte, error) {
	url, err := c.buildURL("simple", map[string]string{
		"i": input,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("Content-Type") != "image/gif" {
		if resp.StatusCode != 200 {
			return nil, ErrInvalidInput
		}
		return nil, ErrInvalidAppID
	}
	return ioutil.ReadAll(resp.Body)
}
