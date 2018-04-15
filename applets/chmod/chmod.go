package chmod

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shirou/toybox/common"
)

const binaryName = "chmod"

var helpFlag bool
var recursiveFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("chmod [-R] mode file...")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&recursiveFlag, "R", false, "Recursively change file mode bits. For each file operand that names a directory, chmod shall change the file mode bits of the directory and all files in the file hierarchy below it.")

	return ret
}

func Main(args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	as := flagSet.Args()
	if flagSet.NArg() < 2 || helpFlag {
		flagSet.Usage()
		return nil
	}

	mode, err := common.ConvertFileModeStr(as[0])
	if err != nil {
		return err
	}

	if recursiveFlag {
		return recurseChmod(as[1:], mode)
	}
	for _, path := range as[1:] {
		path = os.ExpandEnv(path)
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
	}
	return nil
}

func recurseChmod(paths []string, mode os.FileMode) error {
	for _, root := range paths {
		root = os.ExpandEnv(root)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if err := os.Chmod(path, mode); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
