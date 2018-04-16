package main

import (
	goash "github.com/shirou/goash"
	"github.com/shirou/toybox/applets/basename"
	"github.com/shirou/toybox/applets/cat"
	"github.com/shirou/toybox/applets/chgrp"
	"github.com/shirou/toybox/applets/chmod"
	"github.com/shirou/toybox/applets/chown"
	"github.com/shirou/toybox/applets/cksum"
	"github.com/shirou/toybox/applets/cmp"
	"github.com/shirou/toybox/applets/echo"
	"github.com/shirou/toybox/applets/false"
	"github.com/shirou/toybox/applets/ls"
	"github.com/shirou/toybox/applets/mkdir"
	"github.com/shirou/toybox/applets/mv"
	"github.com/shirou/toybox/applets/true"
	"github.com/shirou/toybox/applets/which"
)

var Applets map[string]Applet

type Applet func([]string) error

func init() {
	Applets = map[string]Applet{
		"basename": basename.Main,
		"cat":      cat.Main,
		"chgrp":    chgrp.Main,
		"chown":    chown.Main,
		"chmod":    chmod.Main,
		"cksum":    cksum.Main,
		"cmp":      cmp.Main,
		"echo":     echo.Main,
		"false":    false.Main,
		"ls":       ls.Main,
		"mkdir":    mkdir.Main,
		"mv":       mv.Main,
		"true":     true.Main,
		"which":    which.Main,

		"sh":    goash.Main,
		"ash":   goash.Main,
		"shell": goash.Main,

		"--install": InstallMain,
		"--help":    UsageMain,
	}
}
