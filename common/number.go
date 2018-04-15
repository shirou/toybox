package common

import "strconv"

func IsNumber(s string) (int, bool) {
	if i, err := strconv.Atoi(s); err == nil {
		return i, true
	}
	return 0, false
}
