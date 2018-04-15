package mkdir

import (
	"flag"
	"fmt"
	"os"
)

const binaryName = "mkdir"

var parentFlag bool
var helpFlag bool
var modeFlag string

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("mkdir [OPTION]... DIRECTORY...")
		ret.PrintDefaults()
	}

	ret.BoolVar(&parentFlag, "p", false, "no error if existing, make parent directories as needed")
	ret.BoolVar(&helpFlag, "help", false, "show this message")

	return ret
}

func Main(args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 1 || helpFlag {
		flagSet.Usage()
		return nil
	}

	return mkdir(flagSet.Args(), parentFlag)
}

func mkdir(dirs []string, parent bool) error {
	for _, dir := range dirs {
		if parentFlag {
			err := os.MkdirAll(dir, 0755)
			if !os.IsExist(err) {
				return err
			}
		} else {
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
