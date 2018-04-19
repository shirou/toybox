package chgrp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const binaryName = "chgrp"

var helpFlag bool
var recursiveFlag bool

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("chgrp -R group file...")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")
	ret.BoolVar(&recursiveFlag, "R", false, "Recursively change file group IDs. For each file operand that names a directory, chgrp shall change the group of the directory and all files in the file hierarchy below it. Unless a -H, -L, or -P option is specified, it is unspecified which of these options will be used as the default")

	return ret
}

func Main(stdout io.Writer, args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	as := flagSet.Args()
	if flagSet.NArg() < 2 || helpFlag {
		flagSet.Usage()
		return nil
	}

	gid, err := lookupGid(as[0])
	if err != nil {
		return err
	}

	if recursiveFlag {
		return recurseChgrp(as[1:], gid)
	}
	for _, path := range as[1:] {
		path = os.ExpandEnv(path)
		if err := chgrp(path, gid); err != nil {
			return err
		}
	}
	return nil
}

func recurseChgrp(paths []string, gid int) error {
	for _, root := range paths {
		root = os.ExpandEnv(root)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if err := chgrp(path, gid); err != nil {
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

func chgrp(path string, gid int) error {
	// Not Portable
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return err
	}

	return os.Chown(path, int(st.Uid), gid)
}

func lookupGid(src string) (int, error) {
	group, err := user.LookupGroupId(src)
	if err != nil {
		group, err = user.LookupGroup(src)
		if err != nil {
			return 0, err
		}
	}
	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return 0, err
	}

	return gid, nil
}
