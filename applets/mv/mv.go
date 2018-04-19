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
	helpFlag         bool
	forceFlag        bool
	interactiveFlag  bool
	notOverWriteFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("mv [-fin] SOURCE DEST")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.interactiveFlag, "i", false, "Interactive, prompt before overwrite")
	ret.BoolVar(&opt.forceFlag, "f", false, "Don't prompt before overwriting")
	ret.BoolVar(&opt.notOverWriteFlag, "n", false, "Don't overwrite an existing file")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return mv(flagSet.Args(), opt)
}

func mv(files []string, opt *Option) error {
	src := files[0]
	dst := files[1]
	if !opt.forceFlag && common.FileExists(dst) {
		return os.ErrExist
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}
