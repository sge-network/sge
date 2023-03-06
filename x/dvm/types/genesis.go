package types

import (
	"fmt"

	"github.com/spf13/cast"
)

// DefaultIndex is the default  global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	pubkeyChangeProposalCount := len(gs.ActivePubkeysChangeProposals) +
		len(gs.FinishedPubkeysChangeProposals)
	if cast.ToUint64(pubkeyChangeProposalCount) != gs.ProposalStats.PubkeysChangeCount {
		return fmt.Errorf("sum of active and settled public keys change proposal must be equal to statistics %d <> %d",
			pubkeyChangeProposalCount, gs.ProposalStats.PubkeysChangeCount)
	}

	for _, active := range gs.ActivePubkeysChangeProposals {
		for _, finished := range gs.FinishedPubkeysChangeProposals {
			if active.Id == finished.Proposal.Id {
				return fmt.Errorf("proposal with id is present in both active and finished list %d", &active.Id)
			}
		}
	}

	return gs.Params.Validate()
}
