package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/market interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgAdd{}, "market/Add")
	legacy.RegisterAminoMsg(cdc, &MsgResolve{}, "market/Resolve")
	legacy.RegisterAminoMsg(cdc, &MsgUpdate{}, "market/Update")
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil))
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAdd{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgResolve{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdate{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// amino is the legacy amino codec
	amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
