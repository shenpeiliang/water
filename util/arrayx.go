package util

// 是否在切片中
func InSlice[U comparable](haystack []U, needle U) bool {
	for _, c := range haystack {
		if c == needle {
			return true
		}
	}

	return false
}
