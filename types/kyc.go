package types

// validate checks whether the kyc data is valid for a particular bettor
// if the kyc is required
func (payload KycDataPayload) Validate(address string) bool {
	// ignore is true means that kyc check should be ignored
	if payload.Ignore {
		return true
	}

	if payload.Approved && payload.ID == address {
		return true
	}

	return false
}
