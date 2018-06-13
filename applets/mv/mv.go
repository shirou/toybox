package mv

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shirou/toybox/common"
)

const binaryName = "mv"

type Option struct {
	help         bool
	force        bool
	interactive  bool
	notOverWrite bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("mv [-fin] SOURCE DEST")
		ret.PrintDefaults()
	}
	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")
	ret.BoolVar(&opt.interactive, "i", false, "Interactive, prompt before overwrite")
	ret.BoolVar(&opt.force, "f", false, "Don't prompt before overwriting")
	ret.BoolVar(&opt.notOverWrite, "n", false, "Don't overwrite an existing file")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || opt.help {
		flagSet.Usage()
		return nil
	}

	return mv(stdout, flagSet.Args(), opt)
}

func mv(w io.Writer, files []string, opt *Option) error {
	src := files[0]
	dst := files[1]
	if !opt.force && common.FileExists(dst) {
		return os.ErrExist
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}
