package strutil

// AppendUnique appends unique strings from elems to slice, avoiding duplicates.
func AppendUnique(slice []string, elems ...string) []string {
	uniqueMap := make(map[string]bool)
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

// Contains checks if a slice contains a specified string.
func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// Equal checks if a string is in a slice of strings.
func Equal(v string, elems ...string) bool {
	return Contains(elems, v)
}

// Remove deletes an item from a slice of strings.
func Remove(slice []string, item string) []string {
	for i, a := range slice {
		if a == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// DeleteEmpty removes empty strings from a slice.
func DeleteEmpty(list []string) []string {
	var r []string
	for _, str := range list {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// Unique returns a slice with only unique strings.
func Unique(items []string) []string {
	uniqueMap := make(map[string]bool)
	for _, item := range items {
		if _, ok := uniqueMap[item]; !ok {
			uniqueMap[item] = true
		}
	}
	uniqueItems := make([]string, 0, len(uniqueMap))
	for item := range uniqueMap {
		uniqueItems = append(uniqueItems, item)
	}
	return uniqueItems
}
