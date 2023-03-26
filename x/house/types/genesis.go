package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		DepositList:    []Deposit{},
		WithdrawalList: []Withdrawal{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, w := range gs.WithdrawalList {
		_, err := sdk.AccAddressFromBech32(w.Address)
		if err != nil {
			return fmt.Errorf("invalid withdrawal address %s", w.Address)
		}

		found := false
		for _, d := range gs.DepositList {
			_, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid deposit address %s", d.Creator)
			}

			if w.Address == d.Creator &&
				w.MarketUID == d.MarketUID &&
				w.ParticipationIndex == d.ParticipationIndex {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("the deposit for the depositor address %s, "+
				"market uid %s and participation index %d not found for the withdrawal",
				w.Address,
				w.MarketUID,
				w.ParticipationIndex)
		}
	}

	// TODO: extend validations for market existence
	// and etc.

	return gs.Params.Validate()
}
