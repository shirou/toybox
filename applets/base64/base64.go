package base64

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/shirou/toybox/common"
)

const binaryName = "base64"

type Option struct {
	helpFlag   bool
	decodeFlag bool
	ignoreFlag bool
	wrap       int
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("base64 -c [-d delim] list [file...]")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.decodeFlag, "d", false, "decode")
	ret.IntVar(&opt.wrap, "w", 76, "wrap length")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) (err error) {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	as := flagSet.Args()
	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	var f *os.File
	if len(as) == 0 || as[0] == "-" {
		f = os.Stdin
		if opt.decodeFlag {
			return base64Decode(stdout, f, opt)
		}
		return base64Encode(stdout, f, opt)
	}

	for _, path := range flagSet.Args() {
		f, err = os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()

		if opt.decodeFlag {
			if err := base64Decode(stdout, f, opt); err != nil {
				return err
			}
		} else {
			if err := base64Encode(stdout, f, opt); err != nil {
				return err
			}
		}
	}
	return nil
}

func base64Encode(w io.Writer, f io.Reader, opt *Option) error {
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	base := base64.StdEncoding.EncodeToString(d)
	fmt.Fprintln(w, common.WrapString(base, opt.wrap, "\n"))
	return nil
}
func base64Decode(w io.Writer, f io.Reader, opt *Option) error {
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	base, err := base64.StdEncoding.DecodeString(string(d))
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(base))
	return nil
}
