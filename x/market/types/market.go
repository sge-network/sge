package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMarket(
	uid, creator string,
	startTS, endTS uint64,
	odds []*Odds,
	betConstraints *MarketBetConstraints,
	meta string,
	bookUID string,
	srContributionForHouse sdk.Int,
	status MarketStatus,
) Market {
	return Market{
		UID:                    uid,
		Creator:                creator,
		StartTS:                startTS,
		EndTS:                  endTS,
		Odds:                   odds,
		BetConstraints:         betConstraints,
		Meta:                   meta,
		BookUID:                bookUID,
		SrContributionForHouse: srContributionForHouse,
		Status:                 status,
	}
}
