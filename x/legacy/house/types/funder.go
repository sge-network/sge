package types

type HouseFeeCollectorFunder struct{}

func (HouseFeeCollectorFunder) GetModuleAcc() string {
	return houseFeeCollector
}
