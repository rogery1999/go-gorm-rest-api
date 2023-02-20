package utils

func SplitSlice[T comparable](slice []T, start uint, end uint) []T {
	lastIndex := uint(len(slice) - 1)

	if start > lastIndex {
		return make([]T, 0)
	}

	if end > (lastIndex + 1) {
		return slice[start:]
	}

	return slice[start:end]
}
