package types

import (
	fmt "fmt"

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

	// if there are even number of public keys, there is a chance that proposal votes count get equal
	if count%2 == 0 {
		return fmt.Errorf("public keys count should be an odd number, its count is %d", count)
	}

	return nil
}
