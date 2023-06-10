package types

type OrderBookLiquidityFunder struct{}

func (ol OrderBookLiquidityFunder) GetModuleAcc() string {
	return orderBookLiquidityPool
}
