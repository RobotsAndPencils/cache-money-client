package cache_test

import (
	"net/http"
	"net/http/httptest"
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
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMissingToken(t *testing.T) {
	_, err := cache.NewClient("", endpoint)
	if err != cache.ErrInvalidToken {
		t.Fatalf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestMissingEndpoint(t *testing.T) {
	_, err := cache.NewClient(token, "")
	if err != cache.ErrInvalidEndpoint {
		t.Fatalf("Expected ErrInvalidEndpoint, got %v", err)
	}
}

func TestInvalidEndpoint(t *testing.T) {
	_, err := cache.NewClient(token, "localhost")
	if err != cache.ErrInvalidEndpoint {
		t.Fatalf("Expected ErrInvalidEndpoint, got %v", err)
	}
}

func TestCheck404(t *testing.T) {
	const key = "1234"

	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Errorf("expected HEAD request, got %q", r.Method)
		}
		if r.URL.Path != "/"+key {
			t.Errorf("expected request path %q, got %q", "/"+key, r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth != token {
			t.Errorf("expected Authorization header %q, got %q", token, auth)
		}
		w.WriteHeader(404)
		called = true
	}))
	defer ts.Close()

	client, err := cache.NewClient(token, ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exists, err := client.Check(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected check to return false, got true")
	}
	if !called {
		t.Error("expected test server to be called")
	}
}

func TestCheck204(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := cache.NewClient(token, ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exists, err := client.Check("1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected check to return true, got false")
	}
}
