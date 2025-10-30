package main

// filter filters a slice based on a predicate function.
// It returns a new slice containing only the elements for which the predicate returns true.
// Parameters:
//   - s: The slice to filter.
//   - predicate: The function to test each element of the slice.
//
// Returns:
//   - A new slice containing the filtered elements.
func filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(s)) // Pre-allocate for efficiency
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
