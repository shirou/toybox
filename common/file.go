package common

import (
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return true
	}
	return false
}
