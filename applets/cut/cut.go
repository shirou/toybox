package cut

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shirou/toybox/common"
)

const binaryName = "cut"

type Option struct {
	helpFlag  bool
	bytePos   string
	charPos   string
	fieldPos  string
	delimiter string
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("cut -c [-d delim] list [file...]")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.StringVar(&opt.bytePos, "b", "", "Cut based on a list of bytes.")
	ret.StringVar(&opt.charPos, "c", "", "Cut based on a list of characters.")
	ret.StringVar(&opt.fieldPos, "f", "", "Cut based on a list of fields.")
	ret.StringVar(&opt.delimiter, "d", "\t", "Set the field delimiter to the character delim. The default is the <tab>.")

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
	if len(as) == 0 {
		f = os.Stdin
		return cut(stdout, f, opt)
	}

	for _, path := range flagSet.Args() {
		f, err = os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := cut(stdout, f, opt); err != nil {
			return err
		}
	}

	return nil
}

func cut(w io.Writer, f io.Reader, opt *Option) error {
	target, r, err := parsePos(opt)
	if err != nil {
		return err
	}

	reader := bufio.NewReaderSize(f, 4096)
	for {
		var out []string
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		switch target {
		case "b":
			out = buildByte(line, r)
		case "c":
			out = buildString(string(line), r)
		case "f":
			out = buildField(strings.Split(string(line), opt.delimiter), r)
		}
		fmt.Fprintln(w, strings.Join(out, ""))
	}
	return nil
}

func buildByte(line []byte, r []int) []string {
	out := make([]string, 0, len(r))
	for _, t := range r {
		if t > len(line) {
			return out
		}
		t := t - 1 // start from 1
		out = append(out, string(line[t]))
	}
	return out
}

func buildString(line string, r []int) []string {
	out := make([]string, 0, len(r))
	for _, t := range r {
		if t > len(line) {
			return out
		}
		t := t - 1 // start from 1
		out = append(out, string(line[t]))
	}
	return out
}

func buildField(line []string, r []int) []string {
	out := make([]string, 0, len(r))
	for _, t := range r {
		if t > len(line) {
			return out
		}
		t := t - 1 // start from 1
		out = append(out, line[t])
	}
	return out
}

func parsePos(opt *Option) (string, []int, error) {
	var target string
	var s string

	if opt.bytePos != "" {
		target = "b"
		s = opt.bytePos
	}
	if opt.charPos != "" {
		if len(target) != 0 {
			return "", nil, errors.New("only one type of list may be specified")
		}
		target = "c"
		s = opt.charPos
	}
	if opt.fieldPos != "" {
		if len(target) != 0 {
			return "", nil, errors.New("only one type of list may be specified")
		}
		target = "f"
		s = opt.fieldPos
	}
	r, err := parsePosStr(s)

	return target, r, err
}

// ex: 3
// ex: 1-
// ex: -10
// ex: 1-10
func parsePosStr(s string) ([]int, error) {
	var r []int
	// at first, split by ","
	for _, item := range strings.Split(s, ",") {
		if num, ok := common.IsNumber(item); ok {
			if num == 0 {
				return r, errors.New("positions are numbered from 1")
			}
			if num > 0 { // avoid -11 as minus 11
				r = append(r, num)
				continue
			}
		}
		if !strings.Contains(item, "-") {
			return r, errors.New("failed to parse range")
		}
		n, err := parseRange(item)
		if err != nil {
			return r, err
		}
		r = append(r, n...)
	}

	return common.UniqSortInt(r), nil
}

func parseRange(src string) ([]int, error) {
	var low, high int
	split := strings.Split(strings.TrimSpace(src), "-")
	var r []int

	if src[0] == '-' { // ex: -10
		low = 1
		num, err := parseNum(split[1])
		if err != nil {
			return r, err
		}
		high = num
	} else if src[len(src)-1] == '-' { // ex: 10-
		high = 999 // <TODO> is this correct?
		num, err := parseNum(split[0])
		if err != nil {
			return r, err
		}
		low = num
	} else if len(split) == 2 { // 10-13
		num, err := parseNum(split[0])
		if err != nil {
			return r, err
		}
		low = num

		num, err = parseNum(split[1])
		if err != nil {
			return r, err
		}
		high = num
	} else {
		return r, errors.New("failed to parse range, invalid")
	}
	if high < low {
		return r, errors.New("high end of range less than low end")
	}

	return seq(low, high), nil
}

func seq(low, high int) []int {
	r := make([]int, 0, high-low+1)
	for i := low; i <= high; i++ {
		r = append(r, i)
	}
	return r
}

func parseNum(s string) (int, error) {
	num, ok := common.IsNumber(s)
	if !ok {
		return 0, errors.New("failed to parse range, parseNum")
	}
	if num == 0 {
		return 0, errors.New("positions are numbered from 1")
	}
	return num, nil
}
