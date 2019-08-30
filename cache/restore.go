package cache

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Restore a path to the cache
func Restore(client *Client, key, path string) error {
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("cache-%v-*.zip", key))
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

	zr, err := zip.OpenReader(tmpfile.Name())
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, zf := range zr.File {
		// logger.Printf("  writing %v", zf.Name)
		err := extract(zf, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func extract(zf *zip.File, path string) error {
	file := filepath.Join(path, zf.Name)
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}

	r, err := zf.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(file)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}
