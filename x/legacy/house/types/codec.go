package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// RegisterLegacyAminoCodec registers the necessary x/house interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgDeposit{}, "house/Deposit")
	legacy.RegisterAminoMsg(cdc, &MsgWithdraw{}, "house/Withdraw")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "house/MsgUpdateParams")
}

// RegisterInterfaces registers the x/house interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeposit{},
		&MsgWithdraw{},
		&MsgUpdateParams{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		&DepositAuthorization{},
		&WithdrawAuthorization{},
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
