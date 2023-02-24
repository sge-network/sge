package types

// NewParticipationBetPair creates a new Participation bet pair object
//
//nolint:interfacer
func NewParticipationBetPair(
	bookID, betUID string, participationIndex uint64,
) ParticipationBetPair {
	return ParticipationBetPair{
		BookID:             bookID,
		ParticipationIndex: participationIndex,
		BetUID:             betUID,
	}
}
