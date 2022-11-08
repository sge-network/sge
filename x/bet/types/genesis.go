package types

import (
	"fmt"

	"github.com/sge-network/sge/utils"
)

// DefaultUID is the default  global uid
const DefaultUID uint64 = 1

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		BetList: []Bet{},
		Params:  DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated uid in bet
	betUIDMap := make(map[string]struct{})

	for _, elem := range gs.BetList {
		uid := string(utils.StrBytes(elem.UID))
		if _, ok := betUIDMap[uid]; ok {
			return fmt.Errorf("duplicated uid for bet")
		}
		betUIDMap[uid] = struct{}{}
	}

	return gs.Params.Validate()
}
