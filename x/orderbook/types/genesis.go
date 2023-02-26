package types

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                    DefaultParams(),
		BookList:                  []OrderBook{},
		BookParticipationList:     []BookParticipation{},
		BookExposureList:          []BookOddsExposure{},
		ParticipationExposureList: []ParticipationExposure{},
		Stats:                     OrderBookStats{ResolvedUnsettled: []string{}},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO: extend the validations for all genesis state lists

	return gs.Params.Validate()
}
