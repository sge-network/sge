package keeper

import "github.com/sge-network/sge/x/bet/types"

// kycValidation checks whether the kyc data is valid for a particular bettor
// if the kyc is required
func kycValidation(address string, kycPayload types.KycDataPayload) bool {
	if kycPayload.KycApproved && kycPayload.KycID == address {
		return true
	}

	return false
}
