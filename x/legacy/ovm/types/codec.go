package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/ovm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgSubmitPubkeysChangeProposalRequest{}, "ovm/SubmitPubkeysChangeProposal")
	legacy.RegisterAminoMsg(cdc, &MsgVotePubkeysChangeRequest{}, "ovm/VotePubkeysChange")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "ovm/MsgUpdateParams")
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitPubkeysChangeProposalRequest{},
		&MsgVotePubkeysChangeRequest{},
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// Amino is the legacy Amino codec
var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	// cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
