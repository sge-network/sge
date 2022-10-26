package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers module codec to the app codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddEvent{}, "sportevent/AddEvent", nil)
	cdc.RegisterConcrete(&MsgResolveEvent{}, "sportevent/ResolveEvent", nil)
	cdc.RegisterConcrete(&MsgUpdateEvent{}, "sportevent/UpdateEvent", nil)
	// this line is used by starport scaffolding # 2
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil))
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddEvent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgResolveEvent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateEvent{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// Amino is the legacy aminto codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
