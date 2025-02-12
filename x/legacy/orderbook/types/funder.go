package types

type OrderBookLiquidityFunder struct{}

func (OrderBookLiquidityFunder) GetModuleAcc() string {
	return orderBookLiquidityPool
}
