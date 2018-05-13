package ls

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const binaryName = "ls"

type Option struct {
	helpFlag       bool
	allFlag        bool
	allExcludeFlag bool
	followLinkFlag bool
	xoFlag         bool
	longFlag       bool
	humanFlag      bool
	oneColumnFlag  bool
}

type Directory struct {
	path    string
	Entries []Entry `json:"entries"`
}

const (
	TypeDirectory = "directory"
	TypeRegular   = "regular"
	TypeSymLink   = "symlink"
	TypeHardLink  = "hardlink"
	TypeNamedPipe = "namedpipe"
)

type Entry struct {
	Name       string `json:"name"`
	Mode       string `json:"mode"`
	ModeOctal  int    `json:"mode_octal"`
	User       string `json:"user"`
	Group      string `json:"group"`
	Uid        uint32 `json:"uid"`
	Gid        uint32 `json:"gid"`
	Type       string `json:"type"`
	Size       int64  `json:"size"`
	ModifyTime int64  `json:"modify-time"`
	ModTime    time.Time
}

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

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("ls [-1AaCxdLHRFplinshrSXvctu] [-w WIDTH] [FILE]...")
		ret.PrintDefaults()
	}

	var opt Option
	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.xoFlag, "xo", false, "output libxo compatible json format")
	ret.BoolVar(&opt.allFlag, "a", false, "all")
	ret.BoolVar(&opt.allExcludeFlag, "A", false, "all but exclude . and ..")
	ret.BoolVar(&opt.followLinkFlag, "L", false, "follow symlink")
	ret.BoolVar(&opt.longFlag, "l", false, "long")
	ret.BoolVar(&opt.humanFlag, "h", false, "humanize")
	ret.BoolVar(&opt.oneColumnFlag, "1", false, "one")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return ls(stdout, flagSet.Args(), opt)
}

func ls(w io.Writer, paths []string, opt *Option) error {
	dirs, err := gather(paths, opt)
	if err != nil {
		return err
	}

	return output(w, dirs, opt)
}

func gather(paths []string, opt *Option) ([]Directory, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	ret := make([]Directory, 0)
	for _, path := range paths {
		cur, err := os.Lstat(filepath.Join(path, "."))
		if err != nil {
			return nil, err
		}
		files := []os.FileInfo{cur}
		if path != "/" {
			par, err := os.Lstat(filepath.Join(path, ".."))
			if err != nil {
				return nil, err
			}
			files = append(files, par)
		}

		ff, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		files = append(files, ff...)

		dir := Directory{
			path:    path,
			Entries: make([]Entry, 0, len(files)),
		}
		for _, fi := range files {
			if skip(fi, opt) {
				continue
			}

			var type_ string
			mode := fi.Mode()
			if fi.IsDir() {
				type_ = TypeDirectory
			} else if mode&os.ModeSymlink != 0 {
				type_ = TypeSymLink
				if opt.followLinkFlag {
					fi, err = followLink(path, fi)
					if err != nil {
						return nil, err
					}
					mode = fi.Mode()

				}
			} else if mode&os.ModeNamedPipe != 0 {
				type_ = TypeNamedPipe
			}

			e := Entry{
				Name:       fi.Name(),
				Mode:       mode.String(),
				Size:       fi.Size(),
				Type:       type_,
				ModifyTime: fi.ModTime().Unix(),
				ModTime:    fi.ModTime(),
			}

			addUser(&e)

			dir.Entries = append(dir.Entries, e)
		}
		ret = append(ret, dir)
	}

	return ret, nil
}

func skip(fi os.FileInfo, opt *Option) bool {
	if strings.HasPrefix(fi.Name(), ".") {
		if opt.allFlag {
			return false
		}
		if opt.allExcludeFlag {
			if fi.Name() == "." || fi.Name() == ".." {
				return true
			}
			return false
		}
		return true
	}
	return false
}

func followLink(dir string, fi os.FileInfo) (os.FileInfo, error) {
	path, err := os.Readlink(filepath.Join(dir, fi.Name()))
	if err != nil {
		return nil, err
	}
	return os.Lstat(path)
}
