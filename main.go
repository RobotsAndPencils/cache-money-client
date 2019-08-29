package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RobotsAndPencils/cache-money-client/cache"
)

var logger = log.New(os.Stderr, "", 0)

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
		logger.Printf("warn: %v\n", err)
	}
}

func helpText() {
	logger.Println("cache [store|restore]")
	logger.Println("\nRequired environment variables:")
	logger.Println("  CACHE_KEY={{identifier for the cache}}")
	logger.Println("  CACHE_PATH={{file path to cache}}")
	logger.Println("  TOKEN={{authorization token}}")
	logger.Println("  ENDPOINT={{server url}}")
	logger.Println()
}

func failOnError(err error) {
	if err != nil {
		helpText()
		logger.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
