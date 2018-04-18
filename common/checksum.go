package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
)

func Checksum(t string, fp io.Reader) (string, error) {
	var m hash.Hash
	switch t {
	case "md5":
		m = md5.New()
	case "sha1":
		m = sha1.New()
	case "sha512":
		m = sha512.New()
	case "sha256":
		m = sha256.New()
	case "sha224":
		m = sha256.New224()
	case "sha384":
		m = sha512.New384()
	default:
		return "", fmt.Errorf("unknown type: %s\n", t)
	}

	if _, err := io.Copy(m, fp); err != nil {
		return "", fmt.Errorf("%ssum: %s\n", t, err.Error())
	}

	return fmt.Sprintf("%x", m.Sum(nil)), nil
}
