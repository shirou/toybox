package chown

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const binaryName = "chown"

var helpFlag bool
var recursiveFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("chown -R group file...")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&recursiveFlag, "R", false, "Recursively change file group IDs. For each file operand that names a directory, chown shall change the group of the directory and all files in the file hierarchy below it. Unless a -H, -L, or -P option is specified, it is unspecified which of these options will be used as the default")

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

	uid, err := lookupUid(as[0])
	if err != nil {
		return err
	}

	if recursiveFlag {
		return recurseChown(as[1:], uid)
	}
	for _, path := range as[1:] {
		path = os.ExpandEnv(path)
		if err := chown(path, uid); err != nil {
			return err
		}
	}
	return nil
}

func recurseChown(paths []string, uid int) error {
	for _, root := range paths {
		root = os.ExpandEnv(root)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if err := chown(path, uid); err != nil {
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

func chown(path string, uid int) error {
	// Not Portable
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return err
	}

	return os.Chown(path, uid, int(st.Gid))
}

func lookupUid(src string) (int, error) {
	u, err := user.LookupId(src)
	if err != nil {
		u, err = user.Lookup(src)
		if err != nil {
			return 0, err
		}
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, err
	}

	return uid, nil
}
