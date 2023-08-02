package subaccount

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler initialize a new sdk.handler instance for registered messages
func NewHandler() sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, fmt.Errorf("legacy handler not supported")
	}
}
