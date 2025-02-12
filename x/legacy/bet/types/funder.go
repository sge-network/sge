package types

type BetFeeCollectorFunder struct{}

func (BetFeeCollectorFunder) GetModuleAcc() string {
	return betFeeCollector
}
