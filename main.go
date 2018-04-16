package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func Usage() {
	fmt.Println("toybox -- A minimalistic toolbox, but just a toy")
}
func UsageMain(_ []string) error {
	Usage()
	return nil
}

func main() {
	callname := filepath.Base(os.Args[0])
	args := os.Args[1:]
	applet, ok := Applets[callname]
	if !ok {
		if callname == "toybox" && len(os.Args) > 1 {
			applet, ok = Applets[os.Args[1]]
			if !ok {
				Usage()
				os.Exit(1)
			}
			args = os.Args[2:]
		} else {
			Usage()
			os.Exit(1)
		}
	}

	if err := applet(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
