package keeper

import (
	"context"
	"fmt"

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
	return k.verifyTicketWithKeyUnmarshal(goCtx, ticketStr, clm, "")
}

// verifyTicketWithKeyUnmarshal verifies the ticket using the provided publiv key first, then if the token was verified, it unmarshal the data of ticket into clm.
func (k Keeper) verifyTicketWithKeyUnmarshal(goCtx context.Context, ticketStr string, clm interface{}, pubKey string) error {
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

	// get key vault from module state
	keyVault, found := k.GetKeyVault(ctx)
	if !found {
		return types.ErrNoPublicKeysFound
	}

	pubKeys := keyVault.PublicKeys

	if pubKey != "" {
		// check if the provided pubkey is registered or not
		isRegistered := false
		for _, registereedPubKey := range pubKeys {
			if registereedPubKey == pubKey {
				isRegistered = true
			}
		}

		// pubkey is not registered so it is invalid
		if !isRegistered {
			return fmt.Errorf("the provided public key is not registered in the blockchain store: %s", pubKey)
		}

		// replace the registered pubkeys with the provided key
		pubKeys = []string{pubKey}
	}

	// validate ticket by the keys with provided pubkey
	err = ticket.Verify(pubKeys...)
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
