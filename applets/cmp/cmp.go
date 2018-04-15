package cmp

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shirou/toybox/common"
)

const binaryName = "cmp"

var helpFlag bool
var silenceFlag bool
var linenumberFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("cmp [-l|-s] file1 file2")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&linenumberFlag, "l", false, "Write the byte number (decimal) and the differing bytes (octal) for each difference.")
	ret.BoolVar(&silenceFlag, "s", false, "Write nothing to standard output or standard error when files differ")

	return ret
}

func Main(args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || helpFlag {
		flagSet.Usage()
		return nil
	}

	return cmp(flagSet.Args(), linenumberFlag, silenceFlag)
}

func cmp(files []string, lf, silence bool) error {
	src := files[0]
	dst := files[1]

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	srcR := bufio.NewReaderSize(sf, 4096)

	df, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer df.Close()
	dstR := bufio.NewReaderSize(df, 4096)

	line := 0
	for {
		line += 1
		srcLine, _, srcErr := srcR.ReadLine()
		dstLine, _, dstErr := dstR.ReadLine()

		if srcErr == io.EOF && dstErr == io.EOF {
			break
		} else if srcErr != nil || dstErr != nil {
			return srcErr
		}
		if c, ok := common.CompareSlice(srcLine, dstLine); !ok {
			if silence {
				return nil
			}
			return fmt.Errorf("%s %s differ: char %d line %d\n", src, dst, c, line)
		}
	}
	return nil
}
