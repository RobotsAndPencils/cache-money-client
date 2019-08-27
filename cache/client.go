package cache

import (
	"errors"
	"io"
	"net/url"
)

// Client for cache money server API
type Client struct {
	token    string
	endpoint *url.URL
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
		endpoint: u,
	}, nil
}

// Check if key exists in the cache
func (c *Client) Check(key string) (bool, error) {
	return false, nil
}

// Upload data to the cache
func (c *Client) Upload(key string, r io.Reader) error {
	return nil
}

// Download data from the cache
func (c *Client) Download(key string, w io.Writer) error {
	return nil
}
