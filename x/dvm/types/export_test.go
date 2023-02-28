package types

// JwtTestToken is a wrapper object for the jwtToken, It is being used
// to export unexported methods of the token
type JwtTestToken = jwtTicket

// NewTestJwtToken create new jwt token object
//
//nolint:revive
func NewTestJwtToken(header, payload string, signature string) *JwtTestToken {
	return &jwtTicket{
		header:    header,
		payload:   payload,
		signature: signature,
	}
}

// VerifyWithKey verifies the test token with key
func (t *JwtTestToken) VerifyWithKey(key string) (bool, error) {
	return t.verifyWithKey(key)
}
