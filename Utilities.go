package GolangCPUParts
package GolangCPUParts

// FindIndex returns the index of the first occurrence of `value` in the slice,
// or -1 if the value is not found.
func FindIndex[T comparable](array []T, value T) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1 // Not found
}

func RemoveSlice[T any](array []T, start, end int) []T {
	// Validate indices
	if start < 0 || end > len(array) || start > end {
		panic("invalid start or end indices") // Handle invalid slices (safe for demonstration)
	}

	// Combine the elements before `start` and after `end`
	return append(array[:start], array[end:]...)
}