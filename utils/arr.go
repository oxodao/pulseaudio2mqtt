package utils

func Contains[T comparable](arr []T, elt T) bool {
	for _, v := range arr {
		if v == elt {
			return true
		}
	}

	return false
}
