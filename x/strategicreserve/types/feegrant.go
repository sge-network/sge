package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkfeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
)

func NewFeeGrant(creator, grantee string, blocktime int64) FeeGrant {
	return FeeGrant{
		Creator:        creator,
		Grantee:        grantee,
		GrantBlockTime: blocktime,
	}
}

// DefaultFeeGrantAllowance is the default allowance of the sr pool
func DefaultFeeGrantAllowance(ctx sdk.Context) sdkfeegrant.FeeAllowanceI {
	expireTime := ctx.BlockTime().Add(DefaultAllowanceExpiration * time.Minute)
	return &sdkfeegrant.BasicAllowance{
		Expiration: &expireTime,
	}
}
