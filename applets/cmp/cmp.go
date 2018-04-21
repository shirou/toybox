package cmp

import (
	"bufio"
	"flag"
	"fmt"
	"io"

	"github.com/shirou/toybox/common"
)

const binaryName = "cmp"

type Option struct {
	helpFlag       bool
	silenceFlag    bool
	linenumberFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("cmp [-l|-s] file1 file2")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.linenumberFlag, "l", false, "Write the byte number (decimal) and the differing bytes (octal) for each difference.")
	ret.BoolVar(&opt.silenceFlag, "s", false, "Write nothing to standard output or standard error when files differ")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return cmp(flagSet.Args(), opt.linenumberFlag, opt.silenceFlag)
}

func cmp(files []string, lf, silence bool) error {
	sf, df, err := common.OpenTwoFiles(files)
	if err != nil {
		return err
	}
	defer sf.Close()
	defer df.Close()
	srcR := bufio.NewReaderSize(sf, 4096)
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
			return fmt.Errorf("%s %s differ: char %d line %d\n", files[0], files[1], c, line)
		}
	}
	return nil
}
