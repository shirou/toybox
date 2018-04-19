package rmdir

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const binaryName = "rmdir"

type Option struct {
	helpFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("rmdir dir...")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return rmdir(flagSet.Args())
}

func rmdir(dirs []string) error {
	for _, dir := range dirs {
		dir = os.ExpandEnv(dir)
		si, err := os.Lstat(dir)
		if err != nil {
			return err
		}
		if !si.IsDir() {
			return fmt.Errorf("%s is not directory", dir)
		}
		if err := os.Remove(dir); err != nil {
			return err
		}
	}

	return nil
}
