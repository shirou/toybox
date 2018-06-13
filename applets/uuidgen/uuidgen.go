package uuidgen

import (
	"flag"
	"fmt"
	"io"
)

const binaryName = "uuidgen"

type Option struct {
	help bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("uuidgen")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.help {
		flagSet.Usage()
		return nil
	}

	return uuidgen(stdout)
}
