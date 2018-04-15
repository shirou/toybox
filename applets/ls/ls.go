package ls

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const binaryName = "ls"

var helpFlag bool
var forceFlag bool
var interactiveFlag bool
var notOverWriteFlag bool

/*
   -1      One column output
   -a      Include entries which start with .
   -A      Like -a, but exclude . and ..
   -x      List by lines
   -d      List directory entries instead of contents
   -L      Follow symlinks
   -H      Follow symlinks on command line
   -R      Recurse
   -p      Append / to dir entries
   -F      Append indicator (one of =@|) to entries
   -l      Long listing format
   -i      List inode numbers
   -n      List numeric UIDs and GIDs instead of names
   -s      List allocated blocks
   -lc     List ctime
   -lu     List atime
   --full-time     List full date and time
   -h      Human readable sizes (1K 243M 2G)
   --group-directories-first
   -S      Sort by size
   -X      Sort by extension
   -v      Sort by version
   -t      Sort by mtime
   -tc     Sort by ctime
   -tu     Sort by atime
   -r      Reverse sort order
   -w N    Format N columns wide
   --color[={always,never,auto}]   Control coloring
*/

func NewFlagSet() *flag.FlagSet {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("ls [-1AaCxdLHRFplinshrSXvctu] [-w WIDTH] [FILE]...")
		ret.PrintDefaults()
	}

	ret.BoolVar(&helpFlag, "help", false, "show this message")

	return ret
}

func Main(args []string) error {
	flagSet := NewFlagSet()
	flagSet.Parse(args)

	if helpFlag {
		flagSet.Usage()
		return nil
	}

	return ls(flagSet.Args())
}

func ls(paths []string) error {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	files, _ := ioutil.ReadDir(paths[0])
	return output(os.Stdout, files)
}

func output(w io.Writer, files []os.FileInfo) error {
	for _, f := range files {
		fmt.Fprintln(w, f.Name())
	}

	return nil
}
