package cache

import "os"

// Store a path to the cache
func Store(client *Client, key, path string) error {
	// check if key exists
	exists, err := client.Check(key)
	if err != nil {
		return err
	}
	if exists {
		// No need to cache it
		return nil
	}

	// TODO: compress path
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return client.Push(key, "application/octet-stream", f)
}
