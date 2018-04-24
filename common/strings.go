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

// StringsHas checks the target string slice contains src or not
func StringsHas(target []string, src string) bool {
	for _, t := range target {
		if strings.TrimSpace(t) == src {
			return true
		}
	}
	return false
}

// StringsContains checks the src in any string of the target string slice
func StringsContains(target []string, src string) bool {
	for _, t := range target {
		if strings.Contains(t, src) {
			return true
		}
	}
	return false
}

// IntContains checks the src in any int of the target int slice.
func IntContains(target []int, src int) bool {
	for _, t := range target {
		if src == t {
			return true
		}
	}
	return false
}
