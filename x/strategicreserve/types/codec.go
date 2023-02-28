package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers the codec for the SR module
func RegisterCodec(cdc *codec.LegacyAmino) {
}

// RegisterInterfaces registers the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// Amino stores a reference to a new Legacy Amino
	Amino = codec.NewLegacyAmino()

	// ModuleCdc stores a reference to a new ProtoCodec
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
