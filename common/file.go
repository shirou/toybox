package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const fileModeRegex = `[ugoa]*([-+=]([rwxXst]*|[ugo]))+|[-+=][0-7]+`

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func FileNotExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
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

func IsSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func IsNamedPipe(info os.FileInfo) bool {
	return info.Mode()&os.ModeNamedPipe == os.ModeNamedPipe
}

func IsHidden(name string) bool {
	if name == "." || name == ".." {
		return false
	}
	return len(name) > 1 && name[0] == '.'
}

func IsDir(name string) (bool, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("%q is not a directory", name)
	}
	return true, nil
}

func OpenTwoFiles(files []string) (sf *os.File, df *os.File, err error) {
	src := os.ExpandEnv(files[0])
	dst := os.ExpandEnv(files[1])

	if src == "-" {
		sf = os.Stdin
	} else {
		sf, err = os.Open(src)
		if err != nil {
			return nil, nil, err
		}
	}
	if dst == "-" {
		df = os.Stdin
	} else {
		df, err = os.Open(dst)
		if err != nil {
			return nil, nil, err
		}
	}
	return sf, df, nil
}

const BackupSuffix = "~"

func Backup(path string) error {
	return os.Rename(path, path+BackupSuffix)
}

// see https://github.com/monochromegane/go-gitignore/

type IgnoreMatcher interface {
	Match(path string, isDir bool) bool
}

type IgnoreMatchers []IgnoreMatcher

func (im IgnoreMatchers) Match(path string, isDir bool) bool {
	for _, ig := range im {
		if ig == nil {
			return false
		}
		if ig.Match(path, isDir) {
			return true
		}
	}
	return false
}

type walkFunc func(path string, info os.FileInfo, depth int, ignores IgnoreMatchers) (IgnoreMatchers, error)

func Walk(path string, info os.FileInfo, depth int, parentIgnores IgnoreMatchers, followed bool, walkFn walkFunc) error {
	ignores, walkError := walkFn(path, info, depth, parentIgnores)
	if walkError != nil {
		if info.IsDir() && walkError == filepath.SkipDir {
			return nil
		}
		return walkError
	}

	if !info.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	depth++
	for _, file := range files {
		Walk(filepath.Join(path, file.Name()), file, depth, ignores, followed, walkFn)
	}
	return nil
}
