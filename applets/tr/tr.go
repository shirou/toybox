package tr

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const binaryName = "tr"

type Option struct {
	help bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("tr [SET1] [SET2]")
		ret.PrintDefaults()
	}
	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	as := flagSet.Args()
	if len(as) != 2 || opt.help {
		flagSet.Usage()
		return nil
	}

	return tr(os.Stdout, os.Stdin, as[0], as[1])
}

func tr(w io.Writer, r io.Reader, set1, set2 string) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		dst, err := replace(s.Text(), set1, set2)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, dst)
	}
	return s.Err()
}

func replace(src, old, new string) (string, error) {
	d := strings.Replace(src, old, new, -1)
	return d, nil
}
