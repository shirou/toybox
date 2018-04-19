package wc

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const binaryName = "wc"

type Option struct {
	helpFlag bool
	charFlag bool
	lineFlag bool
	wordFlag bool
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("wc [-c|-m] [-lw] [file...]")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.charFlag, "c", false, "Write to the standard output the number of bytes in each input file.")
	ret.BoolVar(&opt.lineFlag, "l", false, "Write to the standard output the number of <newline> characters in each input file.")
	ret.BoolVar(&opt.wordFlag, "w", false, "Write to the standard output the number of words in each input file.")

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
		return wc(stdout, "", f, opt)
	}

	for _, path := range flagSet.Args() {
		f, err = os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := wc(stdout, path, f, opt); err != nil {
			return err
		}
	}

	return nil
}

type result struct {
	lines     int
	words     int
	bytes     int
	maxLength int
}

func wc(w io.Writer, path string, f io.Reader, opt *Option) error {
	var ret result
	reader := bufio.NewReaderSize(f, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		ret.lines++
		ret.bytes += len(line) + 1 // +1 means NewLine
		ret.words += len(bytes.Fields(line))
	}
	var out []string
	if opt.charFlag {
		out = append(out, strconv.Itoa(ret.bytes))
	}
	if opt.lineFlag {
		out = append(out, strconv.Itoa(ret.lines))
	}
	if opt.wordFlag {
		out = append(out, strconv.Itoa(ret.words))
	}
	if len(out) == 0 { // no flag set
		out = append(out, strconv.Itoa(ret.lines))
		out = append(out, strconv.Itoa(ret.words))
		out = append(out, strconv.Itoa(ret.bytes))
	}
	out = append(out, path)

	fmt.Fprintln(w, strings.Join(out, " "))

	return nil
}
