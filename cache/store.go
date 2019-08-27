package cache

// Store a path to the cache
func Store(client *Client, key, path string) error {
	// check if key exists
	exists, err := client.Check(key)
	if err != nil {
		// ...
	}
	if exists {
		// No need to cache it
		return nil
	}

	// compress path
	// upload data
	return nil
}
