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
	cdc.RegisterConcrete(&MsgPlaceBet{}, "bet/PlaceBet", nil)
	cdc.RegisterConcrete(&MsgPlaceBetSlip{}, "bet/PlaceBetSlip", nil)
	cdc.RegisterConcrete(&MsgSettleBet{}, "bet/SettleBet", nil)
	cdc.RegisterConcrete(&MsgSettleBetBulk{}, "bet/SettleBetBulk", nil)
	// this line is used by starport scaffolding # 2
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPlaceBet{},
		&MsgPlaceBetSlip{},
		&MsgSettleBet{},
		&MsgSettleBetBulk{},
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
