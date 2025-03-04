package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers module codec to the app codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreate{}, "subaccount/Create")
	legacy.RegisterAminoMsg(cdc, &MsgTopUp{}, "subaccount/TopUp")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawUnlockedBalances{}, "subaccount/Withdraw")
	legacy.RegisterAminoMsg(cdc, &MsgWager{}, "subaccount/BetWager")
	legacy.RegisterAminoMsg(cdc, &MsgHouseDeposit{}, "subaccount/HouseDeposit")
	legacy.RegisterAminoMsg(cdc, &MsgHouseWithdraw{}, "subaccount/HouseWithdraw")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "subaccount/MsgUpdateParams")
}

// RegisterInterfaces registers the module interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreate{},
		&MsgTopUp{},
		&MsgWithdrawUnlockedBalances{},
		&MsgWager{},
		&MsgHouseDeposit{},
		&MsgHouseWithdraw{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// Amino is the legacy Amino codec
var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
