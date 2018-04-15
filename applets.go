package main

import (
	goash "github.com/shirou/goash"
	"github.com/shirou/toybox/applets/echo"
	initialize "github.com/shirou/toybox/applets/initialize_toybox"
	"github.com/shirou/toybox/applets/ls"
	"github.com/shirou/toybox/applets/mkdir"
	"github.com/shirou/toybox/applets/mv"
)

var Applets map[string]Applet = map[string]Applet{
	"mkdir":             mkdir.Main,
	"initialize_toybox": initialize.Main,
	"ls":                ls.Main,
	"echo":              echo.Main,
	"mv":                mv.Main,
	"sh":                goash.Main,
	"ash":               goash.Main,
	"shell":             goash.Main,
}

type Applet func([]string) error
