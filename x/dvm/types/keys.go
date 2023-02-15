package types

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
	// PubKeysListKey is the key of the list of public keys in KV-Store which points to []string.
	PubKeysListKey = []byte{0x00}

	// ProposalStatsKey is the key for proposal statistics.
	ProposalStatsKey = []byte{0x01}

	// PubKeysChangeProposalListActivePrefix is the prefix of active pubkeys change proposal.
	PubKeysChangeProposalListActivePrefix = []byte{0x02}

	// FinishedPubKeysChangeProposalListPrefix is the prefix of pubkeys change proposal that is finished.
	FinishedPubKeysChangeProposalListPrefix = []byte{0x03}
)
