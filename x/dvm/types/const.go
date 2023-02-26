package types

// JWT constants
const (
	// JWTHeaderIndex is the index of the header in the JWT ticket
	JWTHeaderIndex = 0

	// JWTPayloadIndex is the index of the payload in the JWT ticket
	JWTPayloadIndex = 1

	// JWTSeparator is the separator character between JWT ticket parts
	JWTSeparator = "."

	// DefaultTimeWeight is the default weight of the time for JWT ticket expiration
	DefaultTimeWeight = 1

	minPubKeysCount = 5

	maxPubKeysCount = 7
)
