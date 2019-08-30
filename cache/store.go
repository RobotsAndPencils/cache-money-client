package cache

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

const mimeType = "application/zip"

var logger = log.New(os.Stderr, "", 0)

// Store a path to the cache
func Store(client *Client, key, path string) error {
	// check if key exists
	logger.Printf("Checking %v", key)
	exists, err := client.Check(key)
	if err != nil {
		return err
	}
	if exists {
		// No need to cache it
		return nil
	}

	tmpfile := filepath.Join(os.TempDir(), fmt.Sprintf("cache-%v.tar.gz", key))
	defer os.Remove(tmpfile)

	// compress path
	logger.Printf("Compressing %v", path)

	err = archiver.Archive([]string{path}, tmpfile)
	if err != nil {
		return err
	}

	logger.Printf("Pushing %v", key)
	f, err := os.Open(tmpfile)
	if err != nil {
		return err
	}
	defer f.Close()
	return client.Push(key, mimeType, f)
}
