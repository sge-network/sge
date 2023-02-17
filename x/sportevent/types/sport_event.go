package types

func NewSpoerEvent(
	uid, creator string,
	startTS, endTS uint64,
	odds []*Odds,
	betConstraits *EventBetConstraints,
	active bool,
	meta string,
) SportEvent {
	return SportEvent{
		UID:            uid,
		Creator:        creator,
		StartTS:        startTS,
		EndTS:          endTS,
		Odds:           odds,
		BetConstraints: betConstraits,
		Active:         active,
		Meta:           meta,
	}
}
