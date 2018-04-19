package rm

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const binaryName = "rm"

var helpFlag bool
var forceFlag bool
var recurseFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("rm [-r] file...")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&recurseFlag, "R", false, "Copy file hierarchies.")
	ret.BoolVar(&recurseFlag, "r", false, "Copy file hierarchies.")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || helpFlag {
		flagSet.Usage()
		return nil
	}

	return rm(flagSet.Args())
}

func rm(files []string) error {
	for _, path := range files {
		path = os.ExpandEnv(path)
		if recurseFlag {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}
