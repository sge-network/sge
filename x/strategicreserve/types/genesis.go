package types

// this line is used by starport scaffolding # genesis/types/import

// NewGenesisState creates a new GenesisState object
func NewGenesisState(reserver Reserver, params Params) *GenesisState {
	return &GenesisState{
		Reserver: reserver,
		Params:   params,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Reserver: DefaultInitialReserver(),
		Params:   DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return ValidateReserver(gs.Reserver)
}
