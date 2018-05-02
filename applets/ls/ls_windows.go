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
		for _, entry := range dir.entries {
			if opt.longFlag && !opt.humanFlag {
				fmt.Fprintf(w, longFormat, entry.mode, entry.size,
					entry.modTime.Format(longTimeFormat), entry.name)
			} else if opt.longFlag && opt.humanFlag {
				fmt.Fprintf(w, longHumanFormat, entry.mode,
					common.Bytes(uint64(entry.size)),
					entry.modTime.Format(longTimeFormat), entry.name)
			} else {
				fmt.Fprintln(w, entry.name)
			}
		}
	}

	return nil
}

func addUser(entry *Entry) {
	// not implemented
}
