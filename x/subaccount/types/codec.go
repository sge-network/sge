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
