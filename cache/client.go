package cache

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// Client for cache money server API
type Client struct {
	token    string
	endpoint string
}

// Errors
var (
	ErrInvalidToken    = errors.New("TOKEN is required")
	ErrInvalidEndpoint = errors.New("ENDPOINT must be a valid URL")
)

// NewClient for cache money server API
func NewClient(token, endpoint string) (*Client, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}
	if endpoint == "" {
		return nil, ErrInvalidEndpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil || !u.IsAbs() {
		return nil, ErrInvalidEndpoint
	}
	return &Client{
		token:    token,
		endpoint: u.String(),
	}, nil
}

// Check if key exists in the cache server
func (c *Client) Check(key string) (bool, error) {
	URL := c.buildURL(key)
	req, err := http.NewRequest("HEAD", URL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 204:
		return true, nil
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf("%v %v", resp.StatusCode, resp.Status)
	}
}

// Upload data to the cache server
func (c *Client) Upload(key, mimeType string, r io.Reader) error {
	URL := c.buildURL(key)
	req, err := http.NewRequest("PUT", URL, r)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", mimeType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%v %v", resp.StatusCode, resp.Status)
	}
	return nil
}

// Download data from the cache server
func (c *Client) Download(key string, w io.Writer) error {
	URL := c.buildURL(key)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%v %v", resp.StatusCode, resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

func (c *Client) buildURL(key string) string {
	// URL already validated in NewClient so won't error here
	// Parsing URL again to avoid mutating c.endpoint
	u, _ := url.Parse(c.endpoint)
	u.Path = path.Join(u.Path, key)
	return u.String()
}
