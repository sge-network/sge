package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

func NewPool(total sdkmath.Int) Pool {
	return Pool{
		Total:     total,
		Spent:     sdkmath.ZeroInt(),
		Withdrawn: sdkmath.ZeroInt(),
	}
}

func (p *Pool) CheckBalance(toSpend sdkmath.Int) error {
	availablePool := p.AvailableAmount()
	if availablePool.LT(toSpend) {
		return sdkerrors.Wrapf(ErrCampaignPoolBalance, "amount %s, available pool %s", toSpend, availablePool)
	}
	return nil
}

func (p *Pool) Spend(amount sdkmath.Int) {
	p.Spent = p.Spent.Add(amount)
}

func (p *Pool) TopUp(amount sdkmath.Int) {
	p.Total = p.Total.Add(amount)
}

func (p *Pool) Withdraw(amount sdkmath.Int) {
	p.Withdrawn = p.Withdrawn.Add(amount)
}

func (p *Pool) AvailableAmount() sdkmath.Int {
	return p.Total.Sub(p.Withdrawn).Sub(p.Spent)
}
