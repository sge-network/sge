package types

type OrderBookLiquidityFunder struct{}

func (OrderBookLiquidityFunder) GetModuleAcc() string {
	return orderBookLiquidityPool
}

type PriceLockFunder struct{}

func (PriceLockFunder) GetModuleAcc() string {
	return priceLockPool
}
