package dirname

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
)

const binaryName = "dirname"

type Option struct {
	helpFlag bool
	zeroFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("dirname [z] NAME")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.zeroFlag, "z", false, "end each output line with NUL, not newline")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() == 0 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return dirname(stdout, flagSet.Args(), opt)
}

func dirname(w io.Writer, files []string, opt *Option) error {
	f := "%s"
	if !opt.zeroFlag {
		f += "\n"
	}

	for _, file := range files {
		p := filepath.Dir(filepath.Clean(file))
		fmt.Fprintf(w, f, p)
	}
	return nil
}
