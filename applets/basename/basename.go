package basename

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const binaryName = "basename"

type Option struct {
	help    bool
	newLine bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("basename [-n] [string ...]")
		ret.PrintDefaults()
	}
	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")
	ret.BoolVar(&opt.newLine, "n", false, "Do not print the trailing newline character.")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	as := flagSet.Args()
	if len(as) > 2 || opt.help {
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
