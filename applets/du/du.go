package du

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/toybox/common"
)

const binaryName = "du"

var sep = string(os.PathSeparator)

type Option struct {
	help   bool
	follow bool
	hidden bool
	depth  int
	human  bool
	zero   bool
	xo     bool
}

type Dir struct {
	name    string
	entries []Entry
	size    int
}

type Entry struct {
	name string
	size int
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("du [OPTIONS] FILE...")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.help, "help", false, "show this message")
	ret.IntVar(&opt.depth, "d", math.MaxInt16, "show this message")
	ret.BoolVar(&opt.xo, "xo", false, "output libxo compatible json format")
	ret.BoolVar(&opt.zero, "0", false, "end each output line with NUL, not newline")
	ret.BoolVar(&opt.human, "h", false, "humanize")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || opt.help {
		flagSet.Usage()
		return nil
	}

	return du(stdout, flagSet.Args(), opt)
}

func du(w io.Writer, files []string, opt *Option) error {
	for _, file := range files {
		if err := walk(w, file, opt); err != nil {
			return err
		}
	}

	return nil
}

func walk(w io.Writer, root string, opt *Option) error {
	maxDepth := opt.depth
	depth := 0
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	ignores := make(common.IgnoreMatchers, 0)

	format := "%s"
	if !opt.zero {
		format += "\n"
	}

	// dir -> size
	sizeMap := make(map[string]int64)

	common.Walk(root, info, depth, ignores, opt.follow, func(path string, info os.FileInfo, depth int, ignores common.IgnoreMatchers) (common.IgnoreMatchers, error) {
		if info.IsDir() {
			if depth > maxDepth {
				return ignores, filepath.SkipDir
			}

			if !opt.hidden && common.IsHidden(info.Name()) {
				return ignores, filepath.SkipDir
			}

			if ignores.Match(path, true) {
				return ignores, filepath.SkipDir
			}
			return ignores, nil
		}
		if !opt.follow && common.IsSymlink(info) {
			return ignores, nil
		}

		if common.IsNamedPipe(info) {
			return ignores, nil
		}

		if !opt.hidden && common.IsHidden(info.Name()) {
			return ignores, filepath.SkipDir
		}

		if ignores.Match(path, false) {
			return ignores, nil
		}

		dir := filepath.Dir(path)
		addSize(dir, info.Size(), sizeMap)

		return ignores, nil
	})

	for dir, si := range sizeMap {
		if opt.human {
			fmt.Fprintf(w, "%-8s %s\n", common.Bytes(uint64(si)), dir)
		} else {
			fmt.Fprintf(w, "%-10d %s\n", si, dir)
		}
	}

	return nil
}

func addSize(dir string, size int64, sizeMap map[string]int64) {
	var d string
	for _, path := range strings.Split(dir, sep) {
		d = filepath.Join(d, path)
		d = strings.TrimSuffix(d, sep)
		sizeMap[d] += size
	}
}
