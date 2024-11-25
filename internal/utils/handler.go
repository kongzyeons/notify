package utils

func ValueInSlice[T comparable](v T, slice []T) bool {
	for i := range slice {
		if v == slice[i] {
			return true
		}
	}
	return false
}
