package date

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const binaryName = "date"

const (
	RFC3339date    = "2006-01-02"
	RFC3339seconds = "2006-01-02 03:04:05Z07:00"
	RFC3339ns      = "2006-01-02 03:04:05.999999999Z07:00"
	ISO8601date    = "2006-01-02"
	ISO8601hour    = "2006-01-02T15Z0700"
	ISO8601minutes = "2006-01-02T15:04Z0700"
	ISO8601seconds = "2006-01-02T15:04:05Z0700"
	ISO8601ns      = "2006-01-02T15:04:05.999999999Z0700"
)

type Option struct {
	helpFlag      bool
	utcFlag       bool
	referenceFlag bool
	iso8601Flag   string
	rfc2822Flag   bool
	rfc3339Flag   string
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("date")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.utcFlag, "u", false, "Display or set the date in UTC (Coordinated Universal) time.")
	ret.BoolVar(&opt.referenceFlag, "r", false, "Print the date and time of the last modification of filename.")
	ret.StringVar(&opt.iso8601Flag, "I", "", "Display date/time in ISO 8601 format.")
	ret.BoolVar(&opt.rfc2822Flag, "rfc-2822", false, "Display date/time in RFC2822")
	ret.BoolVar(&opt.rfc2822Flag, "R", false, "Display date/time in RFC2822")
	ret.StringVar(&opt.rfc3339Flag, "rfc-3339", "", "Display date/time in RFC3339")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return date(stdout, flagSet.Args(), opt)
}

func date(w io.Writer, args []string, opt *Option) error {
	if opt.referenceFlag {
		return reference(w, args, opt)
	}

	var ti time.Time
	if opt.utcFlag {
		ti = time.Now().UTC()
	} else {
		ti = time.Now()
	}

	return output(w, ti, opt)
}

func reference(w io.Writer, args []string, opt *Option) error {
	if len(args) == 0 {
		return errors.New("options require argument")
	}
	if len(args) > 1 {
		return errors.New("options should one argument")
	}
	path := os.ExpandEnv(args[0])
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if opt.utcFlag {
		return output(w, stat.ModTime().UTC(), opt)
	}
	return output(w, stat.ModTime(), opt)
}

func output(w io.Writer, t time.Time, opt *Option) error {
	switch {
	case opt.rfc2822Flag:
		fmt.Fprintln(w, t.Format(time.RFC1123Z))
	case opt.rfc3339Flag != "" && opt.rfc3339Flag != "date" &&
		opt.rfc3339Flag != "seconds" && opt.rfc3339Flag != "ns":
		return fmt.Errorf(`date: invalid argument '%s' for '--rfc-3339'
Valid arguments are:
  - 'date'
  - 'seconds'
  - 'ns'
Try 'date --help' for more information.`, opt.rfc3339Flag)
	case opt.rfc3339Flag == "date":
		fmt.Fprintln(w, t.Format(RFC3339date))
	case opt.rfc3339Flag == "seconds":
		fmt.Fprintln(w, t.Format(RFC3339seconds))
	case opt.rfc3339Flag == "ns":
		fmt.Fprintln(w, t.Format(RFC3339ns))
	case opt.iso8601Flag == "date":
		fmt.Fprintln(w, t.Format(ISO8601date))
	case opt.iso8601Flag == "hours":
		fmt.Fprintln(w, t.Format(ISO8601hour))
	case opt.iso8601Flag == "minutes":
		fmt.Fprintln(w, t.Format(ISO8601minutes))
	case opt.iso8601Flag == "seconds":
		fmt.Fprintln(w, t.Format(ISO8601seconds))
	case opt.iso8601Flag == "ns":
		fmt.Fprintln(w, t.Format(ISO8601ns))
	default:
		fmt.Fprintln(w, t.Format(time.UnixDate))
	}

	return nil
}
