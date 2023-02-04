package types

// DefaultGenesis returns the default  genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:           DefaultParams(),
		Books:            []OrderBook{},
		Bookparticipants: []BookParticipant{},
	}
}
