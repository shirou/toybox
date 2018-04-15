package initialize_toybox

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var AppletBins = []string{
	"basename",
	"cat",
	"chgrp",
	"chmod",
	"chown",
	"cksum",
	"echo",
	"false",
	"initialize_toybox",
	"ls",
	"mkdir",
	"mv",
	"true",
	"which",

	"sh",
	"shell",
}

var AppletSbins = []string{}

var helpFlag bool
var rootDir string

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet("initialize_toybox", flag.ExitOnError)

	ret.Usage = func() {
		fmt.Printf("initialize_toybox\n")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.StringVar(&rootDir, "s", "/", "install target root")

	return ret
}

func Main(args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if helpFlag {
		flagSet.Usage()
		return nil
	}

	return initialize_toybox(rootDir)
}

func initialize_toybox(root string) error {
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

	toypath := filepath.Join(root, "usr", "sbin", "toybox")
	for _, bin := range AppletBins {
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
