package utils

// StrSliceContains check if a slice of strings contain a certain string
func StrSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// StrBytes returns []byte form of its input string
func StrBytes(p string) []byte {
	return []byte(p)
}

// RemoveDuplicateStrs retuns input array without duplicates
func RemoveDuplicateStrs(strSlice []string) (list []string) {
	keys := make(map[string]bool)

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return
}
