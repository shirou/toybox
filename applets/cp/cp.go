package cp

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shirou/toybox/common"
)

const binaryName = "cp"

var helpFlag bool
var forceFlag bool
var recurseFlag bool
var preserveFlag bool
var preserveLinkFlag bool
var newerFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("cp -R [-fip] source_file... target")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&forceFlag, "f", false, "Don't prompt before overwriting")
	ret.BoolVar(&recurseFlag, "R", false, "Copy file hierarchies.")
	ret.BoolVar(&recurseFlag, "r", false, "Copy file hierarchies.")
	ret.BoolVar(&preserveFlag, "p", false, "Duplicate the following characteristics of each source file in the corresponding destination file")
	ret.BoolVar(&preserveLinkFlag, "P", false, "never follow symbolic links in SOURCE")
	ret.BoolVar(&newerFlag, "u", false, "copy only when the SOURCE file is newer than the destination file")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 2 || helpFlag {
		flagSet.Usage()
		return nil
	}

	return cp(flagSet.Args())
}

func cp(files []string) error {
	dst := files[len(files)-1]
	dst = os.ExpandEnv(dst)

	for _, src := range files[:len(files)-1] {
		src = os.ExpandEnv(src)
		if !forceFlag && common.FileExists(dst) {
			return os.ErrExist
		}
		if err := cpOne(src, dst); err != nil {
			return err
		}

	}
	return nil
}

func cpOne(src, dst string) error {
	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	if si.IsDir() {
		if !recurseFlag {
			return errors.New("source is a directory")
		}
		if err := cpDirectory(src, dst); err != nil {
			return err
		}
	} else {
		if err := cpFile(src, dst); err != nil {
			return err
		}
	}

	return nil
}

func cpDirectory(src, dst string) error {
	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	// ensure dst dir does not already exist
	if _, err := os.Open(dst); !os.IsNotExist(err) {
		return errors.New("dstination already exists")
	}

	// create dst dir
	if err := os.MkdirAll(dst, si.Mode()); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := cpOne(filepath.Join(src, file.Name()),
			filepath.Join(dst, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

func cpFile(src, dst string) error {
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if preserveLinkFlag && !si.Mode().IsRegular() {
		return cpSymlink(src, dst)
	}

	di, err := os.Lstat(dst)
	if !os.IsNotExist(err) {
		if !forceFlag {
			return fmt.Errorf("destination already exists: %s", dst)
		}
		if si.ModTime().After(di.ModTime()) && newerFlag {
			return errors.New("destination is newer then src")
		}
	}

	//open source
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	//create dst
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	// copy
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Chmod(si.Mode()); err != nil {
		return err
	}
	if preserveFlag {
		if err = os.Chtimes(dst, si.ModTime(), si.ModTime()); err != nil {
			return err
		}
	}

	//sync dst to disk
	return out.Sync()
}

func cpSymlink(src, dst string) error {
	linkTarget, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(linkTarget, dst)
}
