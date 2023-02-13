package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
)

// VerifyTicket validates a ticket.
// For JWT see https://datatracker.ietf.org/doc/html/rfc7519
// * exp is required.
func (k Keeper) VerifyTicket(goCtx context.Context, ticket string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	t, err := types.NewJwtTicket(ticket)
	if err != nil {
		return err
	}

	// check the expiration of ticket
	err = t.ValidateExpiry(ctx)
	if err != nil {
		return err
	}

	// get pub keys from KV-Store
	keys, found := k.GetKeyVault(ctx)
	if !found {
		return types.ErrNoPublicKeysFound
	}

	// validate the ticket by the keys
	err = t.Verify(keys.PublicKeys...)
	if err != nil {
		return err
	}

	return nil
}

// VerifyTicketUnmarshal verifies the ticket first, then if the token was verified, it unmarshal the data of ticket into clm.
func (k Keeper) VerifyTicketUnmarshal(goCtx context.Context, ticketStr string, clm interface{}) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// construct new ticket object from string ticket
	ticket, err := types.NewJwtTicket(ticketStr)
	if err != nil {
		return err
	}

	// check the expiration of ticket
	err = ticket.ValidateExpiry(ctx)
	if err != nil {
		return err
	}

	// get pub keys from module state
	keys, found := k.GetKeyVault(ctx)

	if !found {
		return types.ErrNoPublicKeysFound
	}

	// validate ticket by the keys
	err = ticket.Verify(keys.PublicKeys...)
	if err != nil {
		return err
	}

	// unmarshal ticket
	err = ticket.Unmarshal(clm)
	if err != nil {
		return err
	}

	return nil
}
