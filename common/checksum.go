package common

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

func Checksum(t string, fp io.Reader) (string, error) {
	var m hash.Hash
	switch t {
	case "md5":
		m = md5.New()
	case "sha1":
		m = sha1.New()
	case "sha256":
		m = sha256.New()
	case "sha512":
		m = sha512.New()
	default:
		return "", fmt.Errorf("unknown type: %s\n", t)
	}

	if _, err := io.Copy(m, fp); err != nil {
		return "", fmt.Errorf("%ssum: %s\n", t, err.Error())
	}

	return fmt.Sprintf("%x", m.Sum(nil)), nil
}

func checkSumOutput(t string, w io.Writer, r io.Reader, path string, suppress bool) error {
	s, err := Checksum(t, r)
	if err != nil {
		return err
	}
	if suppress {
		return nil
	}
	fmt.Fprintf(w, "%s  %s\n", s, path)

	return nil
}

func CheckSumMain(t string, w io.Writer, flagSet *flag.FlagSet, suppress bool) error {
	var f *os.File
	if flagSet.NArg() == 0 || flagSet.Args()[0] == "-" {
		f = os.Stdin
		return checkSumOutput(t, os.Stdout, f, "-", suppress)
	}
	for _, path := range flagSet.Args() {
		r, err := os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer r.Close()

		if err := checkSumOutput(t, os.Stdout, r, path, suppress); err != nil {
			return err
		}
	}
	return nil
}

func CheckSumCompare(t string, w io.Writer, flagSet *flag.FlagSet, suppress bool) error {
	for _, path := range flagSet.Args() {
		f, err := os.Open(os.ExpandEnv(path))
		if err != nil {
			return err
		}
		defer f.Close()

		reader := bufio.NewReaderSize(f, 4096)

		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}

			d := strings.Split(string(line), "  ")
			if len(d) != 2 {
				return fmt.Errorf("wrong md5sum format")
			}
			r, err := os.Open(d[1])
			if err != nil {
				return err
			}
			defer r.Close()

			s, err := Checksum(t, r)
			if err != nil {
				return err
			}
			if s == d[0] {
				fmt.Fprintf(w, "%s: OK\n", d[1])
			} else {
				fmt.Fprintf(w, "%s: Fail\n", d[1])
			}
		}
	}
	return nil
}
