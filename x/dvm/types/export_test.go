package types

// JwtTestToken is a wrapper object for the jwtToken, It is being used
// to export unexported methods of the token
type JwtTestToken = jwtTicket

// NewTestJwtToken create new jwt token object
//
//nolint:revive
func NewTestJwtToken(header, payload string, signatures []string) *JwtTestToken {
	return &jwtTicket{
		header:     header,
		payload:    payload,
		signatures: signatures,
	}
}

// VerifyWithKey verifies the test token with key
func (t *JwtTestToken) VerifyWithKey(key string) (bool, error) {
	return t.verifyWithKey(key)
}
