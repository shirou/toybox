package cat

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const binaryName = "cat"

var helpFlag bool
var unBufferFlag bool

// <TODO> unBufferFlag is not supported

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("cat [-u] [file...]")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&unBufferFlag, "u", false, "Write bytes from the input file to the standard output without delay as each is read.")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if helpFlag {
		flagSet.Usage()
		return nil
	}

	return cat(stdout, flagSet.Args())
}

func cat(w io.Writer, files []string) error {
	for _, path := range files {
		f, err := os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = io.Copy(w, f); err != nil {
			return err
		}
	}
	return nil
}
