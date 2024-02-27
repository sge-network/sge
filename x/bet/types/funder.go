package types

type BetFeeCollectorFunder struct{}

func (BetFeeCollectorFunder) GetModuleAcc() string {
	return betFeeCollector
}

type PriceLockFeeCollector struct{}

func (PriceLockFeeCollector) GetModuleAcc() string {
	return priceLockFeeCollector
}

type PriceLockFunder struct{}

func (PriceLockFunder) GetModuleAcc() string {
	return priceLockPool
}
