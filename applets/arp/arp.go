package arp

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const binaryName = "arp"

type Option struct {
	helpFlag bool
	iname    string
}

type Arp struct {
	addr   string
	hwType string
	flag   string
	hwAddr string
	mask   string
	iname  string
}

type ArpTable map[string]Arp

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("arp [-i]")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.StringVar(&opt.iname, "i", "", "interface")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	table, err := arp(opt)
	if err != nil {
		return err
	}
	if err := output(stdout, table, opt); err != nil {
		return err
	}

	return nil
}

func arp(opt *Option) (ArpTable, error) {
	f, err := os.Open("/proc/net/arp")

	if err != nil {
		return nil, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	s.Scan() // skip the field descriptions

	table := make(map[string]Arp)

	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		a := Arp{
			addr:   fields[0],
			hwType: fields[1],
			flag:   fields[2],
			hwAddr: fields[3],
			mask:   fields[4],
			iname:  fields[5],
		}
		table[a.iname] = a
	}

	return table, nil
}

func output(w io.Writer, table ArpTable, opt *Option) error {
	fmt.Fprintf(w, "%-16s %-17s %s\n", "address", "HW Address", "interface")
	for _, arp := range table {
		if opt.iname != "" && opt.iname != arp.iname {
			continue
		}
		fmt.Fprintf(w, "%-16s %s %s\n", arp.addr, arp.hwAddr, arp.iname)
	}

	return nil
}
