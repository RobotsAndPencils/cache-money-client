package cache

import (
	"fmt"
	"os"
)

// VerifyPath checks that the path exists
func VerifyPath(path string) error {
	if path == "" {
		return fmt.Errorf("CACHE_PATH cannot be blank")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%v does not exist", path)
	}
	return nil
}
