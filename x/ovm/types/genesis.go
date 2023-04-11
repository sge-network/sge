package types

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
	if err := gs.KeyVault.validatePubKeys(); err != nil {
		return err
	}

	return gs.Params.Validate()
}
