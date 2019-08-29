package cache_test

import (
	"testing"

	"github.com/RobotsAndPencils/cache-money-client/cache"
)

func TestSimpleKey(t *testing.T) {
	input := "v1-dependencies"
	output, err := cache.EvaluateKey(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if output != input {
		t.Errorf("expected %v, got %v", input, output)
	}
}

func TestChecksumKey(t *testing.T) {
	input := `v1-{{ checksum "fixtures/go.mod" }}`
	expected := "v1-60c4db30a8f6532d56c9aca1b2aa39d6"

	output, err := cache.EvaluateKey(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if output != expected {
		t.Errorf("expected %v, got %v", expected, output)
	}
}

func TestChecksumFileNotFound(t *testing.T) {
	input := `v1-{{ checksum "fixtures/nofile.lock" }}`
	_, err := cache.EvaluateKey(input)
	if err == nil {
		t.Error("expected error, got none")
	}
}
