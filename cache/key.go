package cache

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"text/template"
)

// EvaluateKey parses the key and evaluates with template functions
func EvaluateKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("CACHE_KEY cannot be blank")
	}

	funcMap := template.FuncMap{
		"checksum": checksum,
	}

	t, err := template.New("EvaluateKey").Funcs(funcMap).Parse(key)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.Execute(&b, nil)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	output := fmt.Sprintf("%x", h.Sum(nil))
	return output, nil
}
