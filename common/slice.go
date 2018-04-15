package common

func CompareSlice(s1, s2 []byte) (int, bool) {
	for i, s := range s1 {
		if s != s2[i] {
			return i, false
		}
	}
	if len(s2) > len(s1) {
		return len(s1) + 1, false
	}

	return 0, true
}
