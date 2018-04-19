package common

import (
	"strings"
)

func WrapString(src string, lim int, delim string) string {
	var buf []string

	for i := 0; i < len(src); i += lim {
		if i+lim < len(src) {
			buf = append(buf, src[i:(i+lim)])
		} else {
			buf = append(buf, src[i:])
		}
	}
	return strings.Join(buf, delim)
}
