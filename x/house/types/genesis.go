package types

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		DepositList:    []Deposit{},
		WithdrawalList: []Withdrawal{},
	}
}
