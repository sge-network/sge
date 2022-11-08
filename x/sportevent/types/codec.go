package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers module codec to the app codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddSportEvent{}, "sportevent/AddSportEvent", nil)
	cdc.RegisterConcrete(&MsgResolveSportEvent{}, "sportevent/ResolveSportEvent", nil)
	cdc.RegisterConcrete(&MsgUpdateSportEvent{}, "sportevent/UpdateSportEvent", nil)
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil))
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddSportEvent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgResolveSportEvent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateSportEvent{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// Amino is the legacy amino codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
