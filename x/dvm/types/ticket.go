package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
)

// Ticket is the Interface of ticket.
type Ticket interface {
	// Unmarshal unmarshals the information of the ticket to the v. v must be a pointer.
	Unmarshal(v interface{}) error

	// Verify verifies the ticket signature with the given public keys. If the ticket is verified by any
	// of the keys, then return nil else return invalid signature error
	Verify(pubKeys ...string) error

	// ValidateExpiry verifies that the ticket is not expired yet.
	ValidateExpiry(ctx sdk.Context) error
}

func (payload *PubkeysChangeProposalPayload) Validate(keys []string, leaderIndex uint32) error {
	count := len(keys)
	if count < minPubKeysCount {
		return fmt.Errorf("total number of pubkeys is %d, this should not be less than %d", count, minPubKeysCount)
	}

	if count > maxPubKeysCount {
		return fmt.Errorf("total number of pubkeys is %d, this should not be more than %d", count, minPubKeysCount)
	}

	if leaderIndex >= cast.ToUint32(count) {
		return fmt.Errorf("leader index is out of range %d", leaderIndex)
	}

	return nil
}

func (payload *ProposalVotePayload) Validate() error {
	if !(payload.Vote == ProposalVote_PROPOSAL_VOTE_YES ||
		payload.Vote == ProposalVote_PROPOSAL_VOTE_NO) {
		return fmt.Errorf("vote should be no or yes %d", payload.Vote)
	}

	return nil
}
