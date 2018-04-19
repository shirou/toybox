package basename

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const binaryName = "basename"

var helpFlag bool
var newLineFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("basename [-n] [string ...]")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&newLineFlag, "n", false, "Do not print the trailing newline character.")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	as := flagSet.Args()
	if len(as) > 2 || helpFlag {
		flagSet.Usage()
		return nil
	}

	str := as[0]
	suffix := ""
	if len(as) == 2 {
		suffix = as[1]
	}

	return basename(os.Stdout, str, suffix)
}

func basename(w io.Writer, str, suffix string) error {
	fmt.Fprintln(w, filepath.Base(os.ExpandEnv(str+suffix)))
	return nil
}
