package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/mint interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
	// left empty because no message is available in this module
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
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
