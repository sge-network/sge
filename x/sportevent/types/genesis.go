package types

import (
	"fmt"

	"github.com/sge-network/sge/utils"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SportEventList: []SportEvent{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated uid in sportEvent
	sportEventUIDMap := make(map[string]struct{})

	for _, elem := range gs.SportEventList {
		uid := string(utils.StrBytes(elem.UID))
		if _, ok := sportEventUIDMap[uid]; ok {
			return fmt.Errorf("duplicated uid for sportEvent")
		}
		sportEventUIDMap[uid] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
