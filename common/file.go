package common

import (
	"fmt"
	"os"
	"strconv"
)

const fileModeRegex = `[ugoa]*([-+=]([rwxXst]*|[ugo]))+|[-+=][0-7]+`

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return true
	}
	return false
}

func ConvertFileModeStr(mode string) (os.FileMode, error) {
	// Octal
	// <TODO> validate
	n, err := strconv.ParseUint(mode, 8, 32)
	if err == nil {
		return os.FileMode(n), err
	}

	// <TODO> string to mode
	return 0, fmt.Errorf("unknown mode: %s", mode)
}

func IsMalformedFileMode(mode os.FileMode) bool {
	return mode != os.ModePerm
}

func IsDir(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("%q is not a directory", name)
	}
	return true, nil
}
