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
