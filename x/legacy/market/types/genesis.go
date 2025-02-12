package types

import (
	"fmt"

	"github.com/sge-network/sge/utils"
)

// DefaultIndex is the default  global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		MarketList: []Market{},
		Stats: MarketStats{
			ResolvedUnsettled: []string{},
		},
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated uid in market
	marketUIDMap := make(map[string]struct{})

	for _, elem := range gs.MarketList {
		uid := string(utils.StrBytes(elem.UID))
		if _, ok := marketUIDMap[uid]; ok {
			return fmt.Errorf("duplicated uid for market")
		}
		marketUIDMap[uid] = struct{}{}
	}

	return gs.Params.Validate()
}
