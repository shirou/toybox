package ln

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shirou/toybox/common"
)

const binaryName = "ln"

type Option struct {
	help     bool
	force    bool
	symbolic bool
	backup   bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("ln [OPTION] TARGET LINK_NAME")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")
	ret.BoolVar(&opt.force, "f", false, "If the target file already exists and is a directory, then remove it so that the link may occur.")
	ret.BoolVar(&opt.symbolic, "s", false, "Create a symbolic link.")
	ret.BoolVar(&opt.backup, "b", false, "make a backup of existing destination file")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || opt.help {
		flagSet.Usage()
		return nil
	}

	return ln(flagSet.Args(), opt)
}

func ln(files []string, opt *Option) error {
	src := os.ExpandEnv(files[0])
	link := os.ExpandEnv(files[1])

	if common.FileNotExists(src) {
		return os.ErrNotExist
	}

	if common.FileExists(link) {
		if opt.backup {
			if err := common.Backup(link); err != nil {
				return err
			}
		} else if !opt.force {
			return os.ErrExist
		} else {
			if err := os.Remove(link); err != nil {
				return err
			}
		}
	}

	if opt.symbolic {
		os.Symlink(src, link)
	} else {
		os.Link(src, link)
	}

	return nil
}
