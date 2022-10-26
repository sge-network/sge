package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers the codec for the SR module
func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
}

// RegisterInterfaces registers the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// Amino stores a reference to a new Legacy Amino
	Amino = codec.NewLegacyAmino()

	// ModuleCdc stores a reference to a new ProtoCodec
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
