package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewParticipantExposure creates a new participant exposure object
//
//nolint:interfacer
func NewParticipantExposure(bookID, oddsID string, exposure, betAmount sdk.Int, pn, round uint64, isFullfilled bool) ParticipantExposure {
	return ParticipantExposure{
		BookID:            bookID,
		OddsID:            oddsID,
		ParticipantNumber: pn,
		Exposure:          exposure,
		BetAmount:         betAmount,
		IsFullfilled:      isFullfilled,
		Round:             round,
	}
}

// MustMarshalParticipantExposure returns the participant exposure bytes. Panics if fails
func MustMarshalParticipantExposure(cdc codec.BinaryCodec, pe ParticipantExposure) []byte {
	return cdc.MustMarshal(&pe)
}

// MustUnmarshalParticipantExposure return the unmarshaled participant exposure from bytes.
// Panics if fails.
func MustUnmarshalParticipantExposure(cdc codec.BinaryCodec, value []byte) ParticipantExposure {
	pe, err := UnmarshalParticipantExposure(cdc, value)
	if err != nil {
		panic(err)
	}

	return pe
}

// return the participant exposure
func UnmarshalParticipantExposure(cdc codec.BinaryCodec, value []byte) (pe ParticipantExposure, err error) {
	err = cdc.Unmarshal(value, &pe)
	return pe, err
}

// String returns a human readable string representation of a ParticipantExposure.
func (pe ParticipantExposure) String() string {
	out, _ := yaml.Marshal(pe)
	return string(out)
}
