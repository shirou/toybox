package main

import (
	"io"

	goash "github.com/shirou/goash"
	"github.com/shirou/toybox/applets/arp"
	"github.com/shirou/toybox/applets/base64"
	"github.com/shirou/toybox/applets/basename"
	"github.com/shirou/toybox/applets/cat"
	"github.com/shirou/toybox/applets/chgrp"
	"github.com/shirou/toybox/applets/chmod"
	"github.com/shirou/toybox/applets/chown"
	"github.com/shirou/toybox/applets/cksum"
	"github.com/shirou/toybox/applets/cmp"
	"github.com/shirou/toybox/applets/cp"
	"github.com/shirou/toybox/applets/cut"
	"github.com/shirou/toybox/applets/date"
	"github.com/shirou/toybox/applets/df"
	"github.com/shirou/toybox/applets/diff"
	"github.com/shirou/toybox/applets/dirname"
	"github.com/shirou/toybox/applets/du"
	"github.com/shirou/toybox/applets/echo"
	"github.com/shirou/toybox/applets/false"
	"github.com/shirou/toybox/applets/head"
	"github.com/shirou/toybox/applets/ln"
	"github.com/shirou/toybox/applets/ls"
	"github.com/shirou/toybox/applets/md5sum"
	"github.com/shirou/toybox/applets/mkdir"
	"github.com/shirou/toybox/applets/mv"
	"github.com/shirou/toybox/applets/rm"
	"github.com/shirou/toybox/applets/rmdir"
	"github.com/shirou/toybox/applets/seq"
	"github.com/shirou/toybox/applets/sha1sum"
	"github.com/shirou/toybox/applets/sha256sum"
	"github.com/shirou/toybox/applets/sha512sum"
	"github.com/shirou/toybox/applets/sleep"
	"github.com/shirou/toybox/applets/sort"
	"github.com/shirou/toybox/applets/tr"
	"github.com/shirou/toybox/applets/true"
	"github.com/shirou/toybox/applets/uniq"
	"github.com/shirou/toybox/applets/uuidgen"
	"github.com/shirou/toybox/applets/wc"
	"github.com/shirou/toybox/applets/wget"
	"github.com/shirou/toybox/applets/which"
	"github.com/shirou/toybox/applets/yes"
)

var Applets map[string]Applet

type Applet func(io.Writer, []string) error

func init() {
	Applets = map[string]Applet{
		"arp":       arp.Main,
		"basename":  basename.Main,
		"base64":    base64.Main,
		"cat":       cat.Main,
		"chgrp":     chgrp.Main,
		"chown":     chown.Main,
		"chmod":     chmod.Main,
		"cksum":     cksum.Main,
		"cmp":       cmp.Main,
		"cp":        cp.Main,
		"cut":       cut.Main,
		"date":      date.Main,
		"df":        df.Main,
		"dirname":   dirname.Main,
		"diff":      diff.Main,
		"du":        du.Main,
		"echo":      echo.Main,
		"false":     false.Main,
		"head":      head.Main,
		"ls":        ls.Main,
		"ln":        ln.Main,
		"mkdir":     mkdir.Main,
		"mv":        mv.Main,
		"md5sum":    md5sum.Main,
		"sha1sum":   sha1sum.Main,
		"sha256sum": sha256sum.Main,
		"sha512sum": sha512sum.Main,
		"sleep":     sleep.Main,
		"seq":       seq.Main,
		"true":      true.Main,
		"uniq":      uniq.Main,
		"uuidgen":   uuidgen.Main,
		"rm":        rm.Main,
		"rmdir":     rmdir.Main,
		"sort":      sort.Main,
		"tr":        tr.Main,
		"yes":       yes.Main,
		"wc":        wc.Main,
		"wget":      wget.Main,
		"which":     which.Main,

		"sh":    goash.Main,
		"ash":   goash.Main,
		"shell": goash.Main,

		"--install": InstallMain,
		"--help":    UsageMain,
	}
}
