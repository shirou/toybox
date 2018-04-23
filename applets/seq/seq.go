package seq

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const binaryName = "seq"

type Option struct {
	helpFlag  bool
	formatStr string
	sep       string
	width     int
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("seq [OPTION] FIRST INCREMENT LAST")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.StringVar(&opt.formatStr, "f", "", "printf format")
	ret.StringVar(&opt.sep, "s", "\n", "separete number")
	ret.IntVar(&opt.width, "w", 0, "width")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	as := flagSet.Args()

	if strings.Contains(strings.Join(args, ""), ".") {
		return seqFloat(stdout, as, opt)
	} else {
		return seqInt(stdout, as, opt)
	}
}

func seqInt(w io.Writer, args []string, opt *Option) (err error) {
	first, increment, last, err := parseArgsInt(args)
	if err != nil {
		return err
	}
	format := opt.formatStr
	if format == "" {
		format = "%d\n"
	}

	if increment > 0 {
		for i := first; i < last; i += increment {
			fmt.Fprintf(w, format, i)
		}
	} else {
		for i := first; i > last; i += increment {
			fmt.Fprintf(w, format, i)
		}
	}

	return nil
}

func parseArgsInt(args []string) (int, int, int, error) {
	first, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, 0, err
	}

	increment := 1
	last := 1

	switch len(args) {
	case 1:
	case 2:
		last, err = strconv.Atoi(args[1])
		if err != nil {
			return 0, 0, 0, err
		}
	case 3:
		increment, err = strconv.Atoi(args[1])
		if err != nil {
			return 0, 0, 0, err
		}
		last, err = strconv.Atoi(args[2])
		if err != nil {
			return 0, 0, 0, err
		}
	default:
		return 0, 0, 0, fmt.Errorf("extra args: %s", strings.Join(args[3:], " "))
	}

	return first, increment, last, nil
}

func seqFloat(w io.Writer, args []string, opt *Option) (err error) {
	first, increment, last, err := parseArgsFloat(args)
	if err != nil {
		return err
	}
	format := opt.formatStr
	if format == "" {
		format = "%.2f\n"
	}

	if increment > 0 {
		for i := first; i < last; i += increment {
			fmt.Fprintf(w, format, i)
		}
	} else {
		for i := first; i > last; i += increment {
			fmt.Fprintf(w, format, i)
		}
	}

	return nil
}

func parseArgsFloat(args []string) (float64, float64, float64, error) {
	first, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		return 0, 0, 0, err
	}

	increment := float64(1)
	last := float64(1)

	switch len(args) {
	case 1:
	case 2:
		last, err = strconv.ParseFloat(args[1], 32)
		if err != nil {
			return 0, 0, 0, err
		}
	case 3:
		increment, err = strconv.ParseFloat(args[1], 32)
		if err != nil {
			return 0, 0, 0, err
		}
		last, err = strconv.ParseFloat(args[2], 32)
		if err != nil {
			return 0, 0, 0, err
		}
	default:
		return 0, 0, 0, fmt.Errorf("extra args: %s", strings.Join(args[3:], " "))
	}

	return first, increment, last, nil
}
