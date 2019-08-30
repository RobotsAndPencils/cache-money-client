package cache

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	// lastDir := ""
	for _, zf := range zr.File {
		// dir := filepath.Dir(zf.Name)
		// if dir != lastDir {
		// 	logger.Printf("  extracting %v", dir)
		// 	lastDir = dir
		// }

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

	logger.Printf("%v %t", file, isSymlink(zf.FileHeader.FileInfo()))

	if isSymlink(zf.FileHeader.FileInfo()) {
		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, r)
		if err != nil {
			return fmt.Errorf("extract symlink %v: %v", file, err)
		}
		return os.Symlink(file, strings.TrimSpace(buf.String()))
	}

	w, err := os.Create(file)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}
