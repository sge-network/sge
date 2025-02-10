package app

import (
	sdkmath "cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
)

func init() {
	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	config.SetAddressVerifier(wasmtypes.VerifyAddressLen())

	err := sdk.RegisterDenom(params.HumanCoinUnit, sdkmath.LegacyOneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(params.BaseCoinUnit, sdkmath.LegacyNewDecWithPrec(1, params.SGEExponent))
	if err != nil {
		panic(err)
	}
	config.Seal()
}
