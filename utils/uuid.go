package utils

import (
	"regexp"
	"strings"
)

// IsValidUID validates if a uid is in correct format
func IsValidUID(uid string) bool {
	if len(uid) == 0 || uid == "" || strings.Contains(uid, " ") ||
		!isValidUUID(uid) {
		return false
	}

	return true
}

// isValidUUID validates the uid
func isValidUUID(uid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uid)
	// _, err := uuid.Parse(uid)
	// return err == nil
}
