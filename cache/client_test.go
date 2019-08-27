package cache_test

import (
	"testing"

	"github.com/RobotsAndPencils/cache-money-client/cache"
)

const (
	token    = "abcd1234"
	endpoint = "https://cache-money.dev/api"
)

func TestValidClient(t *testing.T) {
	_, err := cache.NewClient(token, endpoint)
	if err != nil {
		t.Error("Expected no error")
	}
}

func TestMissingToken(t *testing.T) {
	_, err := cache.NewClient("", endpoint)
	if err != cache.ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestMissingEndpoint(t *testing.T) {
	_, err := cache.NewClient(token, "")
	if err != cache.ErrInvalidEndpoint {
		t.Errorf("Expected ErrInvalidEndpoint, got %v", err)
	}
}

func TestInvalidEndpoint(t *testing.T) {
	_, err := cache.NewClient(token, "localhost")
	if err != cache.ErrInvalidEndpoint {
		t.Errorf("Expected ErrInvalidEndpoint, got %v", err)
	}
}
