package types

type BetFeeCollectorFunder struct{}

func (ol BetFeeCollectorFunder) GetModuleAcc() string {
	return betFeeCollector
}
