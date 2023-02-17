package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewSpoerEvent(
	uid, creator string,
	startTS, endTS uint64,
	odds []*Odds,
	betConstraits *EventBetConstraints,
	active bool,
	meta string,
	bookID string,
	srContributionForHouse sdk.Int,
) SportEvent {
	return SportEvent{
		UID:                    uid,
		Creator:                creator,
		StartTS:                startTS,
		EndTS:                  endTS,
		Odds:                   odds,
		BetConstraints:         betConstraits,
		Active:                 active,
		Meta:                   meta,
		BookId:                 bookID,
		SrContributionForHouse: srContributionForHouse,
	}
}
