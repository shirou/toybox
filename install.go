package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var helpFlag bool
var rootDir string

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet("initialize_toybox", flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("initialize_toybox")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.StringVar(&rootDir, "s", "/", "install target root")

	return ret
}

func InstallMain(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if helpFlag {
		flagSet.Usage()
		return nil
	}

	return install_toybox(rootDir)
}

func install_toybox(root string) error {
	dirs := []string{
		"/usr/bin",
		"/usr/sbin",
	}
	if !filepath.IsAbs(root) {
		pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
		root = filepath.Join(pwd, root)
	}

	for _, dir := range dirs {
		dir = filepath.Join(root, dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	bins := make([]string, 0)
	for k, _ := range Applets {
		if strings.Contains(k, "--") {
			continue
		}

		bins = append(bins, k)
	}

	toypath := filepath.Join(root, "usr", "sbin", "toybox")
	for _, bin := range bins {
		src := filepath.Join(root, "usr", "bin", bin)
		if err := os.Symlink(toypath, src); err != nil {
			if os.IsExist(err) {
				continue
			}
			return err
		}
	}

	return nil
}
