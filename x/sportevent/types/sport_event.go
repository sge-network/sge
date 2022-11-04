package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func InitSportEventConstraints(sportEvent *SportEvent) {
	if sportEvent.BetConstraints == nil {
		sportEvent.BetConstraints = &EventBetConstraints{
			MaxBetCap: DefaultParams().EventMaxBetCap,
			MinAmount: DefaultParams().EventMinBetAmount,
			BetFee:    DefaultParams().EventMinBetFee,
			MaxLoss:   DefaultParams().EventMaxLoss,
			MaxVig:    DefaultParams().EventMaxVig,
			MinVig:    DefaultParams().EventMinVig,
		}
	}

	if sportEvent.BetConstraints.TotalStats == nil {
		sportEvent.BetConstraints.TotalStats = &TotalStats{
			HouseLoss: sdk.ZeroInt(),
			BetAmount: sdk.ZeroInt(),
		}
	}

	if sportEvent.BetConstraints.TotalOddsStats == nil {
		sportEvent.BetConstraints.TotalOddsStats = make(map[string]*TotalOddsStats)
	}
}
