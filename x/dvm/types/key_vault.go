package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
)

// GetLeader sets the current leader of public keys.
func (k *KeyVault) SetLeader(leaderIndex uint32) {
	keys, leaderKey := utils.PopStrAtIndex(k.PublicKeys, leaderIndex)

	// set the leader as first item in the list
	k.PublicKeys = append([]string{leaderKey}, keys...)
}

// GetLeader returns the current leader of public keys.
func (k *KeyVault) GetLeader() string {
	return k.PublicKeys[0]
}

// validatePubKeys validates the public keys stored in the key vault.
func (k *KeyVault) validatePubKeys() error {
	count := len(k.PublicKeys)
	if count < minPubKeysCount {
		return fmt.Errorf("total number of pubkeys is %d, this should not be less than %d", count, minPubKeysCount)
	}

	if count > maxPubKeysCount {
		return fmt.Errorf("total number of pubkeys is %d, this should not be more than %d", count, minPubKeysCount)
	}

	return nil
}

// MajorityCount calculated the minimum count of votes for a proposal to be
// set in order to a proposal pass.
func (k *KeyVault) MajorityCount() int64 {
	count := len(k.PublicKeys)

	majorityVoteCount := sdk.NewDec(int64(count)).
		Mul(minVoteMajorityForDecisionPercentage).
		Ceil().TruncateInt64()

	return majorityVoteCount
}
