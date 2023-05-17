package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterCodec registers module codec to the app codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&DataFeeCollectorFeedProposal{}, "strategicreserve/DataFeeCollectorFeedProposal", nil)
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&DataFeeCollectorFeedProposal{},
	)
}

var (
	// Amino is the legacy amino codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the codec of the module
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
