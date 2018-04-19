package echo

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const binaryName = "echo"

var helpFlag bool
var newLineFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("echo [-n] [string ...]")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&newLineFlag, "n", false, "Do not print the trailing newline character.")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if helpFlag {
		flagSet.Usage()
		return nil
	}

	return echo(stdout, flagSet.Args(), newLineFlag)
}

func echo(w io.Writer, strs []string, newLineFlag bool) error {
	for _, s := range strs {
		if newLineFlag {
			fmt.Fprint(w, os.ExpandEnv(s))
		} else {
			fmt.Fprintln(w, os.ExpandEnv(s))
		}
	}

	return nil
}
