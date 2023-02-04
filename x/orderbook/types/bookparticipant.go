package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewBookParticipant creates a new book participant object
//nolint:interfacer
func NewBookParticipant(bookID string, pAddr sdk.AccAddress, pn uint64, isModuleAccount bool) BookParticipant {
	return BookParticipant{
		BookId:             bookID,
		ParticipantAddress: pAddr.String(),
		ParticipantNumber:  pn,
		IsModuleAccount:    isModuleAccount,
	}
}

// MustMarshalBook returns the orderbook bytes. Panics if fails
func MustMarshalBookParticipant(cdc codec.BinaryCodec, bp BookParticipant) []byte {
	return cdc.MustMarshal(&bp)
}

// MustUnmarshalBookParticipant return the unmarshaled bookparticiapnt from bytes.
// Panics if fails.
func MustUnmarshalBookParticipant(cdc codec.BinaryCodec, value []byte) BookParticipant {
	bp, err := UnmarshalBookParticipant(cdc, value)
	if err != nil {
		panic(err)
	}

	return bp
}

// return the book particiapnt
func UnmarshalBookParticipant(cdc codec.BinaryCodec, value []byte) (bp BookParticipant, err error) {
	err = cdc.Unmarshal(value, &bp)
	return bp, err
}

// String returns a human readable string representation of a BookParticipant.
func (bp BookParticipant) String() string {
	out, _ := yaml.Marshal(bp)
	return string(out)
}
