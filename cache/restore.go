package cache

import "os"

// Restore a path to the cache
func Restore(client *Client, key, path string) error {
	// TODO: decompress to path

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return client.Fetch(key, f)
}
