package main

import (
	"fmt"
	"os"
)

func main() {
	key := os.Getenv("CACHE_KEY")
	path := os.Getenv("CACHE_PATH")
	token := os.Getenv("TOKEN")
	endpoint := os.Getenv("ENDPOINT")

	fmt.Printf("cache_key: %v\n", key)
	fmt.Printf("cache_path: %v\n", path)
	fmt.Printf("token: %v\n", token)
	fmt.Printf("endpoint: %v\n", endpoint)
}
