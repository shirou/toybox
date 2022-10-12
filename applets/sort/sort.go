package sort

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

type Option struct {
	Help bool
}

const binaryName = "sort"

func CreateFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)
	ret.Usage = func() {
		fmt.Println(binaryName)
		ret.PrintDefaults()
	}
	
	var opt Option
	ret.BoolVar(&opt.Help, "help", false, "display this help and exit")
	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := CreateFlagSet()
	flagSet.Parse(args)
	if opt.Help {
		flagSet.Usage()
		return nil
	}

	lines := make([]string, 0)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	sort.Slice(lines, func(i, j int) bool {
		s1 := lines[i]
		s2 := lines[j]
		minLength := len(s1)
		if len(s2) < minLength {
			minLength = len(s2)
		}

		for k := 0; k < minLength; k++ {
			c1 := s1[k]
			c2 := s2[k]
			if c1 < c2 {
				return true
			} else if c1 > c2 {
				return false
			}
		}
		return len(s1) < len(s2)
	})

	for _, line := range lines {
		if _, err := fmt.Fprintln(stdout, line); err != nil {
			return err
		}
	}
	return nil
}
