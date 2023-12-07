package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
)

const (
	// AccountAddressPrefix prefix used for generating account address
	AccountAddressPrefix = "sge"
)

var (
	// AccountPubKeyPrefix used for generating public key
	AccountPubKeyPrefix = AccountAddressPrefix + "pub"
	// ValidatorAddressPrefix used for generating validator address
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	// ValidatorPubKeyPrefix used for generating validator public key
	ValidatorPubKeyPrefix = AccountAddressPrefix + "valoperpub"
	// ConsNodeAddressPrefix used for generating consensus node address
	ConsNodeAddressPrefix = AccountAddressPrefix + "valcons"
	// ConsNodePubKeyPrefix used for generating consensus node public key
	ConsNodePubKeyPrefix = AccountAddressPrefix + "valconspub"
)

// SetConfig sets prefixes configuration
func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)

	err := sdk.RegisterDenom(params.HumanCoinUnit, sdk.OneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(params.BaseCoinUnit, sdk.NewDecWithPrec(1, params.SGEExponent))
	if err != nil {
		panic(err)
	}

	config.Seal()
}
