// +build !windows

package ls

import (
	"fmt"
	"io"
	"os/user"
	"strconv"
	"syscall"

	"github.com/shirou/toybox/common"
)

const longFormat = "%10s %10s %10s %10d %12s %s\n"
const longHumanFormat = "%10s %10s %10s %10s %s %s\n"
const longTimeFormat = "Jan _2 15:04"

func output(w io.Writer, dirs []Directory, opt *Option) error {
	for _, dir := range dirs {
		for _, entry := range dir.entries {
			if opt.longFlag && !opt.humanFlag {
				fmt.Fprintf(w, longFormat, entry.mode,
					entry.user, entry.group,
					entry.size,
					entry.modTime.Format(longTimeFormat), entry.name)
			} else if opt.longFlag && opt.humanFlag {
				fmt.Fprintf(w, longHumanFormat, entry.mode,
					entry.user, entry.group,
					common.Bytes(uint64(entry.size)),
					entry.modTime.Format(longTimeFormat), entry.name)
			} else {
				fmt.Fprintln(w, entry.name)
			}
		}
	}

	return nil
}

var uidCache = make(map[uint32]string)
var gidCache = make(map[uint32]string)

func addUser(entry *Entry) {
	var st syscall.Stat_t
	if err := syscall.Stat(entry.name, &st); err != nil {
		return
	}
	entry.uid = st.Uid
	entry.gid = st.Gid

	uname, ok := uidCache[st.Uid]
	if !ok {
		u, err := user.LookupId(strconv.Itoa(int(st.Uid)))
		if err != nil {
			return
		}
		uname = u.Name
	}
	entry.user = uname

	group, ok := gidCache[st.Gid]
	if !ok {
		g, err := user.LookupGroupId(strconv.Itoa(int(st.Gid)))
		if err != nil {
			return
		}
		group = g.Name
	}
	entry.group = group
}
