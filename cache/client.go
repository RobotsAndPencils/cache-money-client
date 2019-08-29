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

// Check if key exists in the cache
func (c *Client) Check(key string) (bool, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return false, err
	}
	u.Path = path.Join(u.Path, key)
	req, err := http.NewRequest("HEAD", u.String(), nil)
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

// Upload data to the cache
func (c *Client) Upload(key string, r io.Reader) error {
	return nil
}

// Download data from the cache
func (c *Client) Download(key string, w io.Writer) error {
	return nil
}
