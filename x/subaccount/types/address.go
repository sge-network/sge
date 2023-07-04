package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAddressFromSubaccount returns an account address for a subaccount
func NewAddressFromSubaccount(subaccountID uint64) sdk.AccAddress {
	return types.NewModuleAddress(fmt.Sprintf("%s/%d", ModuleName, subaccountID))
}
