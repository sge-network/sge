package utils

// StrBytes returns []byte form of its input string
func StrBytes(p string) []byte {
	return []byte(p)
}

// RemoveDuplicateStrs returns input array without duplicates
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

// RemoveStr removes an item from string slice.
func RemoveStr(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
