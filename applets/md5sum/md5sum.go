package md5sum

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shirou/toybox/common"
)

const binaryName = "md5sum"

type Option struct {
	helpFlag     bool
	compareFlag  bool
	suppressFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("md5sum [-cs] FILE")
		ret.PrintDefaults()
	}
	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.compareFlag, "c", false, "Compare the digest of the file against this string")
	ret.BoolVar(&opt.suppressFlag, "s", false, "Show nothing. only status code")

	return ret, &opt
}

func Main(args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	// if compare flag set, separed function is called
	if opt.compareFlag {
		return compare(os.Stdout, flagSet.Args())
	}

	var f *os.File
	if flagSet.NArg() == 0 || flagSet.Args()[0] == "-" {
		f = os.Stdin
		return md5sum(os.Stdout, f, "-", opt)
	}
	for _, path := range flagSet.Args() {
		f, err := os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := md5sum(os.Stdout, f, path, opt); err != nil {
			return err
		}
	}

	return nil
}

func md5sum(w io.Writer, r io.Reader, path string, opt *Option) error {
	s, err := common.Checksum("md5", r)
	if err != nil {
		return err
	}
	if opt.suppressFlag {
		return nil
	}
	fmt.Fprintf(w, "%s  %s\n", s, path)

	return nil
}

func compare(w io.Writer, files []string) error {
	for _, path := range files {
		f, err := os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()

		reader := bufio.NewReaderSize(f, 4096)

		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}

			d := strings.Split(string(line), "  ")
			if len(d) != 2 {
				return fmt.Errorf("wrong md5sum format")
			}
			r, err := os.Open(d[1])
			if err != nil {
				return err
			}
			defer r.Close()

			s, err := common.Checksum("md5", r)
			if err != nil {
				return err
			}
			if s == d[0] {
				fmt.Fprintf(w, "%s: OK\n", d[1])
			} else {
				fmt.Fprintf(w, "%s: Fail\n", d[1])
			}
		}
	}

	return nil
}
