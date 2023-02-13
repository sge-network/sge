package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/types"
)

const (
	minApprovePublicKeyCount = 3
	maxValidProposalTime     = 30
)

// SetActivePubkeysChangeProposal sets a pubkey list change proposal in the store.
func (k Keeper) SetActivePubkeysChangeProposal(ctx sdk.Context, proposal types.PublicKeysChangeProposal) {
	store := k.getActivePubKeysChangeProposalStore(ctx)
	b := k.cdc.MustMarshal(&proposal)
	store.Set(utils.Uint64ToBytes(proposal.Id), b)
}

// GetActivePubkeysChangeProposal returns a pubkeys change proposat by its id
func (k Keeper) GetActivePubkeysChangeProposal(ctx sdk.Context, id uint64) (val types.PublicKeysChangeProposal, found bool) {
	store := k.getActivePubKeysChangeProposalStore(ctx)

	b := store.Get(utils.Uint64ToBytes(id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllActivePubkeysChangeProposals returns list of all active pubkeys change proposals
func (k Keeper) GetAllActivePubkeysChangeProposals(ctx sdk.Context) (list []types.PublicKeysChangeProposal, err error) {
	store := k.getActivePubKeysChangeProposalStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PublicKeysChangeProposal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// removeActiveProposal removes an active pubkeys change proposal.
func (k Keeper) removeActiveProposal(ctx sdk.Context, id uint64) {
	store := k.getActivePubKeysChangeProposalStore(ctx)
	store.Delete(utils.Uint64ToBytes(id))
}

// SetFinishedPubkeysChangeProposal sets a finished pubkey list change proposal in the store.
func (k Keeper) SetFinishedPubkeysChangeProposal(ctx sdk.Context, finishedProposal types.PublicKeysChangeFinishedProposal) {
	store := k.getFinishedPubKeysChangeProposalStore(ctx)
	b := k.cdc.MustMarshal(&finishedProposal)
	store.Set(utils.Uint64ToBytes(finishedProposal.Proposal.Id), b)
}

// GetFinishedPubkeysChangeProposal returns a finished pubkeys change proposat by its id
func (k Keeper) GetFinishedPubkeysChangeProposal(ctx sdk.Context, id uint64) (val types.PublicKeysChangeFinishedProposal, found bool) {
	store := k.getFinishedPubKeysChangeProposalStore(ctx)

	b := store.Get(utils.Uint64ToBytes(id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllFinishedPubkeysChangeProposals returns list of all finished pubkeys change proposals
func (k Keeper) GetAllFinishedPubkeysChangeProposals(ctx sdk.Context) (list []types.PublicKeysChangeFinishedProposal, err error) {
	store := k.getActivePubKeysChangeProposalStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PublicKeysChangeFinishedProposal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) FinishProposals(ctx sdk.Context) error {
	return k.finishPubkeysChangeProposals(ctx)
}

func (k Keeper) finishPubkeysChangeProposals(ctx sdk.Context) error {
	activeProposals, err := k.GetAllActivePubkeysChangeProposals(ctx)
	if err != nil {
		return err
	}

	if len(activeProposals) == 0 {
		return nil
	}

	for _, proposal := range activeProposals {

		expireTime := ctx.BlockTime().Add(maxValidProposalTime * time.Minute).Unix()
		if proposal.StartTS > expireTime {
			// proposal is expired
			k.finishPubkeysChangeProposal(
				ctx,
				proposal.Id,
				types.ProposalResult_PROPOSAL_RESULT_EXPIRED,
				fmt.Sprintf("current block time is %d, more than %d minutes is passed since start time.", expireTime, maxValidProposalTime),
			)
		}

		if len(proposal.ApprovedBy) > minApprovePublicKeyCount {

			pubKeys, found := k.GetKeyVault(ctx)
			if !found {
				fmt.Printf("there is no publick keys record")
			}
			for _, added := range proposal.Modifications.Additions {
				if err := types.IsValidJwtToken(added); err != nil {
					k.finishPubkeysChangeProposal(
						ctx,
						proposal.Id,
						types.ProposalResult_PROPOSAL_RESULT_REJECTED,
						fmt.Sprintf("public key %s is not a valid jwt token.", added),
					)

					break
				}
				pubKeys.PublicKeys = append(pubKeys.PublicKeys, added)
			}

			for _, deleted := range proposal.Modifications.Deletions {
				pubKeys.PublicKeys = utils.RemoveStr(pubKeys.PublicKeys, deleted)
			}

			k.SetKeyVault(ctx, pubKeys)

			k.finishPubkeysChangeProposal(
				ctx,
				proposal.Id,
				types.ProposalResult_PROPOSAL_RESULT_APPROVED,
				"",
			)
		}
	}

	return nil
}

func (k Keeper) finishPubkeysChangeProposal(
	ctx sdk.Context,
	proposalID uint64,
	result types.ProposalResult,
	resultMetadata string,
) error {
	proposal, found := k.GetActivePubkeysChangeProposal(ctx, proposalID)
	if !found {
		return fmt.Errorf("proposal not found with id %d", proposalID)
	}

	k.SetFinishedPubkeysChangeProposal(ctx,
		types.NewFinishedPublicKeysChangeProposal(
			proposal,
			result,
			resultMetadata,
			ctx.BlockTime().Unix(),
		),
	)
	k.removeActiveProposal(ctx, proposalID)

	return k.finishPubkeysChangeProposals(ctx)
}
