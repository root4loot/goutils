package sliceutil

// AppendUnique appends unique elements from elems to slice, avoiding duplicates.
func AppendUnique[T comparable](slice []T, elems ...T) []T {
	uniqueMap := make(map[T]bool)
	for _, v := range slice {
		uniqueMap[v] = true
	}
	for _, e := range elems {
		if !uniqueMap[e] {
			slice = append(slice, e)
			uniqueMap[e] = true
		}
	}
	return slice
}

// Contains checks if a slice contains a specified element.
func Contains[T comparable](slice []T, item T) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// Remove deletes an item from a slice of elements.
func Remove[T comparable](slice []T, item T) []T {
	for i, a := range slice {
		if a == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// DeleteEmpty removes zero-value elements from a slice.
func DeleteEmpty[T comparable](list []T) []T {
	var r []T
	for _, elem := range list {
		var zeroValue T
		if elem != zeroValue {
			r = append(r, elem)
		}
	}
	return r
}

// Unique returns a slice with only unique elements, preserving the order of the first occurrence.
func Unique[T comparable](items []T) []T {
	uniqueMap := make(map[T]bool)
	var uniqueItems []T
	for _, item := range items {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			uniqueItems = append(uniqueItems, item)
		}
	}
	return uniqueItems
}
