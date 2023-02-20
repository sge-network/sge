package types

// NewParticipationBetPair creates a new Participation bet pair object
//
//nolint:interfacer
func NewParticipationBetPair(
	bookID, betUID string, participationIndex, betID uint64,
) ParticipationBetPair {
	return ParticipationBetPair{
		BookID:             bookID,
		ParticipationIndex: participationIndex,
		BetUID:             betUID,
		BetID:              betID,
	}
}
