package types

import "github.com/sge-network/sge/utils"

// module constants
const (
	// ModuleName defines the module name
	ModuleName = "dvm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dvm"
)

// store prefixes and keys.
var (
	// KeyVaultKey is the key of the list of public keys in KV-Store which points to []string.
	KeyVaultKey = []byte{0x00}

	// ProposalStatsKey is the key for proposal statistics.
	ProposalStatsKey = []byte{0x01}

	// PubKeysChangeProposalListPrefix is the prefix of pubkeys change proposal.
	PubKeysChangeProposalListPrefix = []byte{0x02}
)

// PubkeysChangeProposalPrefix returns prefix of the proposal list.
func PubkeysChangeProposalPrefix(status ProposalStatus) []byte {
	return utils.Int32ToBytes(int32(status))
}

// PubkeysChangeProposalKey returns key of the proposal list.
func PubkeysChangeProposalKey(status ProposalStatus, id uint64) []byte {
	return append(PubkeysChangeProposalPrefix(status), utils.Uint64ToBytes(id)...)
}
