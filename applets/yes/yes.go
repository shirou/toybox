package yes

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

const binaryName = "yes"

type Option struct {
	helpFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("yes [OPTION] strings...")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return yes(stdout, flagSet.Args())
}

func yes(w io.Writer, str []string) error {
	var s string
	if len(str) == 0 {
		s = "yes"
	} else {
		s = strings.Join(str, " ")
	}
	for {
		fmt.Fprintln(w, s)
	}
	return nil
}
