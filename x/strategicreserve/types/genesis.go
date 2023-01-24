package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(reserver Reserver, params Params) *GenesisState {
	return &GenesisState{
		Reserver: reserver,
		Params:   params,
	}
}

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Reserver: DefaultInitialReserver(),
		Params:   DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return ValidateReserver(gs.Reserver)
}
