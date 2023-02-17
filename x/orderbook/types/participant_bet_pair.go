package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// NewParticipantBetPair creates a new participant bet pair object
//
//nolint:interfacer
func NewParticipantBetPair(
	bID, betUID string, participantNumber, betID uint64,
) ParticipantBetPair {
	return ParticipantBetPair{
		BookID:            bID,
		ParticipantNumber: participantNumber,
		BetUID:            betUID,
		BetID:             betID,
	}
}

// MustMarshalParticipantBetPair returns the participant bet pair bytes. Panics if fails
func MustMarshalParticipantBetPair(cdc codec.BinaryCodec, pbp ParticipantBetPair) []byte {
	return cdc.MustMarshal(&pbp)
}

// return the participant fuillfilled bet
func UnmarshalParticipantBetPair(cdc codec.BinaryCodec, value []byte) (pbp ParticipantBetPairResponse, err error) {
	err = cdc.Unmarshal(value, &pbp)
	return pbp, err
}
