package df

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shirou/toybox/common"
)

const binaryName = "df"

type Option struct {
	helpFlag  bool
	allFlag   bool
	humanFlag bool
}

type UsageStat struct {
	Path              string
	Fstype            string
	Total             uint64
	Free              uint64
	Used              uint64
	UsedPercent       float64
	InodesTotal       uint64
	InodesUsed        uint64
	InodesFree        uint64
	InodesUsedPercent float64
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("df [-h]")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.allFlag, "a", false, "includes all fstype")
	ret.BoolVar(&opt.humanFlag, "h", false, "print sizes in powers of 1024")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return df(stdout, flagSet.Args(), opt)
}

const format = "%-16s %10dK %10dK %10dK %3.0f%% %s\n"
const humanFormat = "%-16s %6s %6s %6s %3.0f%% %s\n"

func df(w io.Writer, paths []string, opt *Option) error {
	paths, err := listMountPaths(opt)
	if err != nil {
		return err
	}

	for _, path := range paths {
		ps, err := Usage(path)
		if err != nil {
			if os.IsPermission(err) {
				continue
			}
			return err
		}
		if opt.humanFlag {
			fmt.Fprintf(w, humanFormat, ps.Fstype,
				common.BytesShort(ps.Total/1024),
				common.BytesShort(ps.Used/1024),
				common.BytesShort(ps.Free/1024),
				ps.UsedPercent, ps.Path)
		} else {
			fmt.Fprintf(w, format, ps.Fstype, ps.Total/1024, ps.Used/1024, ps.Free/1024, ps.UsedPercent, ps.Path)
		}
	}

	return nil
}
