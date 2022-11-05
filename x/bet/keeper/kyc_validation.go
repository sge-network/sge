package keeper

import (
	"github.com/sge-network/sge/x/bet/types"
)

// KycValidation checks whether the kyc data is valid for a particular bettor
// if the kyc is required
func KycValidation(address string,
	ticket *types.BetPlacementTicketPayload) bool {

	if ticket.KycData.KycApproved && ticket.KycData.KycId == address {
		return true
	}

	return false
}
