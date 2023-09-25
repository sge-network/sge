package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

func NewPool(total sdkmath.Int) Pool {
	return Pool{
		Total: total,
	}
}

func (p *Pool) CheckBalance(toSpend sdkmath.Int) error {
	availablePool := p.Total.Sub(p.Spent)
	if availablePool.LT(toSpend) {
		return sdkerrors.Wrapf(ErrCampaignPoolBalance, "amount %s, available pool %s", toSpend, availablePool)
	}
	return nil
}
