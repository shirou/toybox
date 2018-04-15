package common

import (
	"bufio"
)

func ReadWholeLine(b *bufio.Reader) (line string, err error) {
	byteline := make([]byte, 0)
	prefix := true
	for prefix {
		var partial []byte
		partial, prefix, err = b.ReadLine()
		if err != nil {
			break
		}
		byteline = append(byteline, partial...)
	}
	return string(byteline), err
}
