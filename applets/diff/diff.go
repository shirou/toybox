package diff

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/shirou/toybox/common"
)

const binaryName = "diff"

type Hunk struct {
	start int
}

type Option struct {
	helpFlag    bool
	unifiedFlag bool
	unifiedNum  int
	contextFlag bool
	contextNum  int
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("diff [-abBdiNqrTstw] [-L LABEL] [-S FILE] [-U LINES] FILE1 FILE2")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.unifiedFlag, "u", true, "Output 3 lines of unified context")
	ret.IntVar(&opt.unifiedNum, "U", 3, "specify lines in unified context")
	ret.BoolVar(&opt.contextFlag, "c", true, "Output NUM (default 3) lines of copied context.")
	ret.IntVar(&opt.contextNum, "C", 3, "specify lines in context")

	return ret, &opt
}

/*
 * Main implements very naive diff algorithm
 * we should use other implementation such as diffmatchpatch
 * https://github.com/sergi/go-diff
 */
func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() != 2 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	if err := diff(stdout, flagSet.Args(), opt); err != nil {
		return err
	}

	return nil
}

func diff(w io.Writer, files []string, opt *Option) error {
	ff, tf, err := common.OpenTwoFiles(files)
	if err != nil {
		return err
	}
	defer ff.Close()
	defer tf.Close()

	from, err := ioutil.ReadAll(ff)
	if err != nil {
		return err
	}
	to, err := ioutil.ReadAll(tf)
	if err != nil {
		return err
	}

	text, err := lineDiff(string(from), string(to), files[0], files[1], opt)
	if err != nil {
		return err
	}

	fmt.Fprint(w, text)

	return nil
}
func lineDiff(from, to, fromFile, toFile string, opt *Option) (string, error) {
	if opt.unifiedFlag {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(from),
			B:        difflib.SplitLines(to),
			FromFile: fromFile,
			ToFile:   toFile,
			Context:  opt.unifiedNum,
		}
		return difflib.GetUnifiedDiffString(diff)
	}
	if opt.contextFlag {
		diff := difflib.ContextDiff{
			A:        difflib.SplitLines(from),
			B:        difflib.SplitLines(to),
			FromFile: fromFile,
			ToFile:   toFile,
			Context:  opt.contextNum,
		}
		return difflib.GetContextDiffString(diff)
	}

	return "", errors.New("no output format specified")
}
