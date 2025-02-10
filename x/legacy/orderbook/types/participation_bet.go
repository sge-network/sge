package types

// NewParticipationBetPair creates a new Participation bet pair object
//
//nolint:interface
func NewParticipationBetPair(
	orderBookUID, betUID string, participationIndex uint64,
) ParticipationBetPair {
	return ParticipationBetPair{
		OrderBookUID:       orderBookUID,
		ParticipationIndex: participationIndex,
		BetUID:             betUID,
	}
}
