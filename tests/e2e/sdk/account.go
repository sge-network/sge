package sdk

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdknetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateAccount(val *sdknetwork.Validator, uid string) sdk.AccAddress {
	// Create new account in the keyring.
	k, _, err := val.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	if err != nil {
		panic(err)
	}

	addr, err := k.GetAddress()
	if err != nil {
		panic(err)
	}

	return addr
}
