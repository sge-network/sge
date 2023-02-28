package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers module codec to the app codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPlaceBet{}, "bet/PlaceBet", nil)
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPlaceBet{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// Amino is the legacy aminto codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
