package main

import (
	"fmt"
	"os"

	"github.com/RobotsAndPencils/cache-money-client/cache"
)

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "store" && os.Args[1] != "restore") {
		failOnError(fmt.Errorf("error: store or restore is required"))
	}

	action := os.Args[1]
	token := os.Getenv("TOKEN")
	endpoint := os.Getenv("ENDPOINT")
	cacheKey := os.Getenv("CACHE_KEY")
	cachePath := os.Getenv("CACHE_PATH")

	client, err := cache.NewClient(token, endpoint)
	failOnError(err)

	err = cache.VerifyPath(cachePath)
	failOnError(err)

	key, err := cache.EvaluateKey(cacheKey)
	failOnError(err)

	switch action {
	case "store":
		err = cache.Store(client, key, cachePath)
	case "restore":
		err = cache.Restore(client, key, cachePath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: %v\n", err)
	}
}

func helpText() {
	fmt.Fprintln(os.Stderr, "cache [store|restore]")
	fmt.Fprintln(os.Stderr, "\nRequired environment variables:")
	fmt.Fprintln(os.Stderr, "  CACHE_KEY={{identifier for the cache}}")
	fmt.Fprintln(os.Stderr, "  CACHE_PATH={{file path to cache}}")
	fmt.Fprintln(os.Stderr, "  TOKEN={{authorization token}}")
	fmt.Fprintln(os.Stderr, "  ENDPOINT={{server url}}")
	fmt.Fprintln(os.Stderr)
}

func failOnError(err error) {
	if err != nil {
		helpText()
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
