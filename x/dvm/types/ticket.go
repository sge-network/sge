package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (payload *PubkeysChangeProposalPayload) Validate(keys []string) error {
	finalCount := len(keys) + len(payload.Additions) - len(payload.Deletions)
	if finalCount < minPubKeysCount {
		return fmt.Errorf("total number of pubkeys is %d, this should not be less than %d", finalCount, minPubKeysCount)
	}

	if finalCount > maxPubKeysCount {
		return fmt.Errorf("total number of pubkeys is %d, this should not be more than %d", finalCount, minPubKeysCount)
	}

	// loop through additions and check if is a valid jwt token
	for _, v := range payload.Additions {
		// check if pem content is a valid ED25516 key
		if err := IsValidJwtToken(v); err != nil {
			return fmt.Errorf("public key %s is not valid jwt token %s", v, err)
		}
	}

	return nil
}
