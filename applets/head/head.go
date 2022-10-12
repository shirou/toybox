package head

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const binaryName = "head"

type Option struct {
	help  bool
	lines int
	bytes int64
	quiet bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("head [-nb] FILES...")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")
	ret.IntVar(&opt.lines, "n", 10, "print the first N lines instead of the first 10; with the leading '-' print all but the last N lines of each file")
	ret.Int64Var(&opt.bytes, "c", 0, "print the first N bytes of each file; with the leading '-', print all but the last N bytes of each file")
	ret.BoolVar(&opt.quiet, "q", false, "not output file name")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || opt.help {
		flagSet.Usage()
		return nil
	}

	var err error
	for _, path := range flagSet.Args() {
		if flagSet.NArg() > 1 && !opt.quiet {
			fmt.Fprintf(stdout, "\n==> %s <==\n", path)
		}
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

		if err := head(stdout, f, opt); err != nil {
			return err
		}
	}

	return nil
}

func head(w io.Writer, f *os.File, opt *Option) error {
	if opt.bytes > 0 {
		return bytesHead(w, f, opt.bytes)
	} else if opt.bytes < 0 {
		return negativeBytesHead(w, f, opt.bytes)
	}
	if opt.lines > 0 {
		return linesHead(w, f, opt.lines)
	} else if opt.lines < 0 {
		return negativeLinesHead(w, f, opt.lines)
	}

	return nil
}

func linesHead(w io.Writer, f *os.File, n int) error {
	s := bufio.NewScanner(f)
	cur := 1
	for s.Scan() {
		if _, err := fmt.Fprintln(w, s.Text()); err != nil {
			return err
		}
		if n <= cur {
			break
		}
		cur += 1
	}
	return s.Err()
}

func negativeLinesHead(w io.Writer, f *os.File, n int) error {
	n *= -1
	lastNLines := make([]string, n)
	pos := 0

	s := bufio.NewScanner(f)
	for count := 0; s.Scan(); count++ {
		if count >= n {
			if _, err := fmt.Fprintln(w, lastNLines[pos]); err != nil {
				return err
			}
		}
		lastNLines[pos] = s.Text()
		pos++
		if pos >= n {
			pos %= n
		}
	}
	return s.Err()
}

func bytesHead(w io.Writer, f *os.File, n int64) error {
	lr := io.LimitReader(f, n)
	if _, err := io.Copy(w, lr); err != nil {
		return err
	}
	return nil
}

func negativeBytesHead(w io.Writer, f *os.File, n int64) error {
	n *= -1
	lastNBytes := make([]byte, n)
	pos := int64(0)

	r := bufio.NewReader(f)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for count := int64(0); ; count++ {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if count >= n {
			if err := bw.WriteByte(lastNBytes[pos]); err != nil {
				return err
			}
		}
		lastNBytes[pos] = b
		pos++
		if pos >= n {
			pos %= n
		}
	}
	return bw.Flush()
}
