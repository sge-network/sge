package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec registers module codec to the app codec
func RegisterCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil))
}

var (
	// Amino is the legacy amino codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
