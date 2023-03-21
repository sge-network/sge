package types

// NewParticipationBetPair creates a new Participation bet pair object
//
//nolint:interfacer
func NewParticipationBetPair(
	bookUID, betUID string, participationIndex uint64,
) ParticipationBetPair {
	return ParticipationBetPair{
		BookUID:            bookUID,
		ParticipationIndex: participationIndex,
		BetUID:             betUID,
	}
}
