package types

type HouseFeeCollectorFunder struct{}

func (ol HouseFeeCollectorFunder) GetModuleAcc() string {
	return houseFeeCollector
}
