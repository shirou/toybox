package common

import (
	"sort"
	"strconv"
)

func IsNumber(s string) (int, bool) {
	if i, err := strconv.Atoi(s); err == nil {
		return i, true
	}
	return 0, false
}

func UniqSortInt(src []int) []int {
	m := make(map[int]struct{})
	for _, s := range src {
		m[s] = struct{}{}
	}
	uniq := make([]int, 0, len(m))
	for i := range m {
		uniq = append(uniq, i)
	}

	sort.Ints(uniq)
	return uniq
}
