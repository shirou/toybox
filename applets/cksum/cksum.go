package cksum

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
)

const binaryName = "cksum"

var helpFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("cksum [file...]")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")

	return ret
}

func Main(args []string) (err error) {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || helpFlag {
		flagSet.Usage()
		return nil
	}

	for _, path := range flagSet.Args() {
		var f *os.File
		if path == "-" {
			f = os.Stdin
		} else {
			f, err = os.Open(os.ExpandEnv(path))
			if err != nil {
				return err
			}
			defer f.Close()
		}

		if err := cksum(os.Stdout, f, crc32.IEEE); err != nil {
			return err
		}
	}
	return nil
}

func cksum(w io.Writer, f *os.File, polynomial uint32) error {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	// <TODO> mismatch
	crc32q := crc32.MakeTable(crc32.IEEE)
	fmt.Fprintf(w, "%d %d %s\n", crc32.Checksum(content, crc32q), len(content), f.Name())
	return nil
}
