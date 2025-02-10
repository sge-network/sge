package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateCampaign{}, "reward/CreateCampaign")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateCampaign{}, "reward/UpdateCampaign")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawFunds{}, "reward/WithdrawFunds")
	legacy.RegisterAminoMsg(cdc, &MsgGrantReward{}, "reward/GrantReward")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "reward/MsgUpdateParams")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCampaign{},
		&MsgUpdateCampaign{},
		&MsgWithdrawFunds{},
		&MsgGrantReward{},
		&MsgUpdateParams{},
	)

	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		&CreateCampaignAuthorization{},
		&UpdateCampaignAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
