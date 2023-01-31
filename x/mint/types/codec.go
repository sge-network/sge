package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers module codec to the app codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	// left empty because no message is available in this module
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// Amino is the legacy aminto codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
