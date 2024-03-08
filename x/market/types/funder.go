package types

type PriceLockFunder struct{}

func (PriceLockFunder) GetModuleAcc() string {
	return priceLockPool
}

type PriceLockFeeCollector struct{}

func (PriceLockFeeCollector) GetModuleAcc() string {
	return priceLockFeeCollector
}
