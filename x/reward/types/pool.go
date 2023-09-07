package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (p *Pool) HasEnoughFund(toSpend sdk.Int) bool {
	return p.Total.Sub(p.Spent).LTE(toSpend)
}
