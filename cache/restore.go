package cache

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mholt/archiver"
)

// Restore a path to the cache
func Restore(client *Client, key, path string) error {
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("cache-%v-*.tar.gz", key))
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	logger.Printf("Fetching %v", key)
	err = client.Fetch(key, tmpfile)
	if err != nil {
		tmpfile.Close()
		return err
	}
	tmpfile.Close()

	//  decompress to path
	logger.Printf("Extracting to %v", path)

	return archiver.Unarchive(tmpfile.Name(), path)
}
