package mv

import (
	"flag"
	"fmt"
	"os"

	"github.com/shirou/toybox/common"
)

const binaryName = "mv"

var helpFlag bool
var forceFlag bool
var interactiveFlag bool
var notOverWriteFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("mv [-fin] SOURCE DEST")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&interactiveFlag, "i", false, "Interactive, prompt before overwrite")
	ret.BoolVar(&forceFlag, "f", false, "Don't prompt before overwriting")
	ret.BoolVar(&notOverWriteFlag, "n", false, "Don't overwrite an existing file")

	return ret
}

func Main(args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || helpFlag {
		flagSet.Usage()
		return nil
	}

	return mv(flagSet.Args())
}

func mv(files []string) error {
	src := files[0]
	dst := files[1]
	if !forceFlag && common.FileExists(dst) {
		return os.ErrExist
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}
