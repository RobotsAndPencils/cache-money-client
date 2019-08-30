package cache

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("cache-%v-*.zip", key))
	if err != nil {
		return err
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	// compress path
	logger.Printf("Compressing %v", path)
	err = compress(path, tmpfile)
	if err != nil {
		return err
	}

	_, err = tmpfile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	logger.Printf("Pushing %v", key)
	return client.Push(key, mimeType, tmpfile)
}

func compress(root string, w io.Writer) error {
	zw := zip.NewWriter(w)
	// lastDir := ""
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// dir := filepath.Dir(rel)
		// if dir != lastDir {
		// 	logger.Printf("  adding %v", dir)
		// 	lastDir = dir
		// }

		zf, err := zw.Create(rel)
		if err != nil {
			return err
		}

		if isSymlink(info) {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("compress symlink %v: %v", path, err)
			}
			_, err = zf.Write([]byte(filepath.ToSlash(linkTarget)))
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(zf, f)
		f.Close()
		return err
	})
	if err != nil {
		return err
	}
	return zw.Close()
}

func isSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0
}
