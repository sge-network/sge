package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sge-network/sge/utils"
)

// SetLeader sets the current leader of public keys.
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
	if count < MinPubKeysCount {
		return fmt.Errorf(
			"total number of pubkeys is %d, this should not be less than %d",
			count,
			MinPubKeysCount,
		)
	}

	if count > MaxPubKeysCount {
		return fmt.Errorf(
			"total number of pubkeys is %d, this should not be more than %d",
			count,
			MaxPubKeysCount,
		)
	}

	for _, pubKey := range k.PublicKeys {
		ed25519Key, err := jwt.ParseEdPublicKeyFromPEM([]byte(pubKey))
		if err != nil {
			return fmt.Errorf("unable to parse public key %s: %v", ed25519Key, err)
		}
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
