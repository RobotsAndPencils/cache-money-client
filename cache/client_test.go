package cache_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		t.Error("test server did not receive a request")
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

func TestUpload(t *testing.T) {
	const key = "1234"
	const content = "These pretzels are making me thirsty."
	const mimeType = "plain/text"

	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT request, got %q", r.Method)
		}
		if r.URL.Path != "/"+key {
			t.Errorf("expected request path %q, got %q", "/"+key, r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth != token {
			t.Errorf("expected Authorization header %q, got %q", token, auth)
		}
		length := r.Header.Get("Content-Length")
		if length != strconv.Itoa(len(content)) {
			t.Errorf("expected Content-Length header %v, got %v", len(content), length)
		}
		mime := r.Header.Get("Content-Type")
		if mime != mimeType {
			t.Errorf("expected Content-Type header %q, got %q", mimeType, mime)
		}
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unexpected error reading body: %v", err)
		}
		if string(b) != content {
			t.Errorf("expected body %q, got %q", content, b)
		}
		w.WriteHeader(200)
		called = true
	}))
	defer ts.Close()

	client, err := cache.NewClient(token, ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = client.Upload(key, mimeType, bytes.NewBufferString(content))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Error("test server did not receive a request")
	}
}

func TestUploadFails(t *testing.T) {
	const key = "1234"
	const content = "These pretzels are making me thirsty."
	const mimeType = "plain/text"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	client, err := cache.NewClient(token, ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = client.Upload(key, mimeType, bytes.NewBufferString(content))
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
