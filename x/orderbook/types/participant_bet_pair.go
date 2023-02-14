package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// NewParticipantBetPair creates a new participant bet pair object
//nolint:interfacer
func NewParticipantBetPair(
	bId, betUuid string, participantNumber, betId uint64,
) ParticipantBetPair {
	return ParticipantBetPair{
		BookId:            bId,
		ParticipantNumber: participantNumber,
		BetUuid:           betUuid,
		BetId:             betId,
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
