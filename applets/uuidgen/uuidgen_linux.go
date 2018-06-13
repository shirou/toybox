package uuidgen

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func uuidgen(w io.Writer) error {
	fp, err := os.Open("/proc/sys/kernel/random/uuid")
	if err != nil {
		return err
	}
	defer fp.Close()
	reader := bufio.NewReaderSize(fp, 38)
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(line))
	return nil
}
