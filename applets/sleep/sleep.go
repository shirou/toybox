package sleep

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const binaryName = "sleep"

var helpFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("sleep number")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 1 || helpFlag {
		flagSet.Usage()
		return nil
	}

	return sleep(flagSet.Arg(0))
}

func sleep(timeStr string) error {
	duration, err := parseTime(timeStr)
	if err != nil {
		return err
	}
	if duration < 0 {
		return errors.New("numer is under 0")
	}
	time.Sleep(duration)

	return nil
}

func parseTime(timeStr string) (time.Duration, error) {
	last := len(timeStr) - 1
	if strings.HasSuffix(timeStr, "s") {
		return parseTimeStr(timeStr[:last], time.Second)
	} else if strings.HasSuffix(timeStr, "m") {
		return parseTimeStr(timeStr[:last], time.Minute)
	} else if strings.HasSuffix(timeStr, "h") {
		return parseTimeStr(timeStr[:last], time.Hour)
	} else if strings.HasSuffix(timeStr, "d") {
		return parseTimeStr(timeStr[:last], time.Hour*24)
	}

	d, err := strconv.ParseFloat(timeStr, 32)
	if err != nil {
		return time.Second, err
	}

	return time.Duration(d) * time.Second, nil
}

func parseTimeStr(timeStr string, duration time.Duration) (time.Duration, error) {
	d, err := strconv.ParseFloat(timeStr, 32)
	if err != nil {
		return time.Second, err
	}
	return time.Duration(d) * duration, nil
}
