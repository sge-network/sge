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
