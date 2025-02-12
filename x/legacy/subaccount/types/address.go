package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAddressFromSubaccount returns an account address for a subaccount
func NewAddressFromSubaccount(subaccountID uint64) sdk.AccAddress {
	return address.Module(types.ModuleName, sdk.Uint64ToBigEndian(subaccountID))
}
