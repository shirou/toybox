package main

import (
	goash "github.com/shirou/goash"
	"github.com/shirou/toybox/applets/basename"
	"github.com/shirou/toybox/applets/cat"
	"github.com/shirou/toybox/applets/chgrp"
	"github.com/shirou/toybox/applets/chmod"
	"github.com/shirou/toybox/applets/chown"
	"github.com/shirou/toybox/applets/echo"
	"github.com/shirou/toybox/applets/false"
	initialize "github.com/shirou/toybox/applets/initialize_toybox"
	"github.com/shirou/toybox/applets/ls"
	"github.com/shirou/toybox/applets/mkdir"
	"github.com/shirou/toybox/applets/mv"
	"github.com/shirou/toybox/applets/true"
	"github.com/shirou/toybox/applets/which"
)

var Applets map[string]Applet = map[string]Applet{
	"basename":          basename.Main,
	"cat":               cat.Main,
	"chgrp":             chgrp.Main,
	"chown":             chown.Main,
	"chmod":             chmod.Main,
	"echo":              echo.Main,
	"false":             false.Main,
	"initialize_toybox": initialize.Main,
	"ls":                ls.Main,
	"mkdir":             mkdir.Main,
	"mv":                mv.Main,
	"true":              true.Main,
	"which":             which.Main,

	"sh":    goash.Main,
	"ash":   goash.Main,
	"shell": goash.Main,
}

type Applet func([]string) error
