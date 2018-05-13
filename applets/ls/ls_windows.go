// +build windows

package ls

import (
	"fmt"
	"io"

	"github.com/shirou/toybox/common"
)

const longFormat = "%10s %10d %12s %s\n"
const longHumanFormat = "%10s %10s %s %s\n"
const longTimeFormat = "Jan _2 15:04"

func output(w io.Writer, dirs []Directory, opt *Option) error {
	for _, dir := range dirs {
		for _, entry := range dir.Entries {
			if opt.longFlag && !opt.humanFlag {
				fmt.Fprintf(w, longFormat, entry.Mode, entry.Size,
					entry.ModTime.Format(longTimeFormat), entry.Name)
			} else if opt.longFlag && opt.humanFlag {
				fmt.Fprintf(w, longHumanFormat, entry.Mode,
					common.Bytes(uint64(entry.Size)),
					entry.ModTime.Format(longTimeFormat), entry.Name)
			} else {
				fmt.Fprintln(w, entry.Name)
			}
		}
	}

	return nil
}

func addUser(entry *Entry) {
	// not implemented
}
