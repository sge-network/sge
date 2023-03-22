package types

import "github.com/sge-network/sge/utils"

func (k *KeyVault) SetLeader(leaderIndex uint32) {
	keys, leaderKey := utils.PopStrAtIndex(k.PublicKeys, leaderIndex)

	// set the leader as first item in the list
	k.PublicKeys = append([]string{leaderKey}, keys...)
}
