package types

type BetFeeCollectorFunder struct{}

func (BetFeeCollectorFunder) GetModuleAcc() string {
	return betFeeCollector
}

type PriceLockFunder struct{}

func (PriceLockFunder) GetModuleAcc() string {
	return priceLockPool
}
