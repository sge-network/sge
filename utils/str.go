package utils

import "strings"

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
		entry = strings.TrimSpace(entry)
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return
}

// PopStrAtIndex pops an string item from string slice by index.
func PopStrAtIndex(s []string, i uint32) ([]string, string) {
	popElem := s[i]
	return append(s[:i], s[i+1:]...), popElem
}
