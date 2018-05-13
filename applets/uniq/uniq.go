package uniq

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const binaryName = "uniq"

type Option struct {
	help       bool
	count      bool
	repeated   bool
	unique     bool
	ignoreCase bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("uniq [-ic] FILE ")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")
	ret.BoolVar(&opt.count, "c", false, "Precede each output line with a count of the number of times the line occurred in the input.")
	ret.BoolVar(&opt.repeated, "d", false, "Suppress the writing of lines that are not repeated in the input.")
	ret.BoolVar(&opt.unique, "u", false, "Suppress the writing of lines that are repeated in the input.")
	ret.BoolVar(&opt.ignoreCase, "i", false, "ignore case")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.help {
		flagSet.Usage()
		return nil
	}

	switch flagSet.NArg() {
	case 0:
		in := os.Stdin
		out := stdout
		return uniq(in, out, opt)
	case 1:
		in, err := os.Open(os.ExpandEnv(flagSet.Arg(0)))
		if err != nil {
			return err
		}
		defer in.Close()
		out := stdout

		return uniq(in, out, opt)
	case 2:
		in, err := os.Open(os.ExpandEnv(flagSet.Arg(0)))
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Open(os.ExpandEnv(flagSet.Arg(1)))
		if err != nil {
			return err
		}
		defer out.Close()
		return uniq(in, out, opt)
	default:
		flagSet.Usage()
		return nil
	}
}

func uniq(in io.Reader, out io.Writer, opt *Option) error {
	s := bufio.NewScanner(in)
	//	var count int
	var dup bool

	// read first line as last line
	s.Scan()
	last := s.Text()

	for s.Scan() {
		if opt.ignoreCase {
			dup = strings.EqualFold(last, s.Text())
		} else {
			dup = last == s.Text()
		}
		if dup {
			if opt.repeated {
				fmt.Fprintln(out, last)
			}
		} else {
			if opt.unique {
				fmt.Fprintln(out, last)
			}
		}
		last = s.Text()
	}

	return nil
}
