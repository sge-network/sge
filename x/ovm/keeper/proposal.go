package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/ovm/types"
)

// SetPubkeysChangeProposal sets a pubkey list change proposal in the store.
func (k Keeper) SetPubkeysChangeProposal(ctx sdk.Context, proposal types.PublicKeysChangeProposal) {
	store := k.getPubKeysChangeProposalStore(ctx)
	b := k.cdc.MustMarshal(&proposal)
	store.Set(types.PubkeysChangeProposalKey(proposal.Status, proposal.Id), b)
}

// GetPubkeysChangeProposal returns a pubkeys change proposat by its id
func (k Keeper) GetPubkeysChangeProposal(
	ctx sdk.Context,
	status types.ProposalStatus,
	id uint64,
) (val types.PublicKeysChangeProposal, found bool) {
	store := k.getPubKeysChangeProposalStore(ctx)

	b := store.Get(types.PubkeysChangeProposalKey(status, id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllPubkeysChangeProposalsByStatus returns list of all pubkeys change proposals
func (k Keeper) GetAllPubkeysChangeProposalsByStatus(
	ctx sdk.Context,
	status types.ProposalStatus,
) (list []types.PublicKeysChangeProposal, err error) {
	store := k.getPubKeysChangeProposalStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.PubkeysChangeProposalPrefix(status))

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

// GetAllPubkeysChangeProposals returns list of all pubkeys change proposals
func (k Keeper) GetAllPubkeysChangeProposals(
	ctx sdk.Context,
) (list []types.PublicKeysChangeProposal, err error) {
	store := k.getPubKeysChangeProposalStore(ctx)
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

// RemoveProposal removes a pubkeys change proposal.
func (k Keeper) RemoveProposal(ctx sdk.Context, status types.ProposalStatus, id uint64) {
	store := k.getPubKeysChangeProposalStore(ctx)
	store.Delete(types.PubkeysChangeProposalKey(status, id))
}

// FinishProposals sets all of active proposals as finished.
func (k Keeper) FinishProposals(ctx sdk.Context) error {
	return k.finishPubkeysChangeProposals(ctx)
}

// finishPubkeysChangeProposals sets all active pubkeys change proposals as finished.
func (k Keeper) finishPubkeysChangeProposals(ctx sdk.Context) error {
	activeProposals, err := k.GetAllPubkeysChangeProposalsByStatus(
		ctx,
		types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
	)
	if err != nil {
		return fmt.Errorf("error while retrieving the active pubkeys change proposals: %s", err)
	}

	if len(activeProposals) == 0 {
		return nil
	}

	blockTime := ctx.BlockTime().Unix()

	keyVault, found := k.GetKeyVault(ctx)
	if !found {
		fmt.Printf("there is no publick keys record")
	}

	for _, proposal := range activeProposals {
		if proposal.IsExpired(blockTime) {
			// proposal is expired
			if err = k.finishPubkeysChangeProposal(
				ctx,
				proposal.Id,
				types.ProposalResult_PROPOSAL_RESULT_EXPIRED,
				fmt.Sprintf("current block time is %d, more than %d minutes is passed since start time.", blockTime, types.MaxValidProposalMinutes),
			); err != nil {
				return fmt.Errorf("error while setting the proposal as expired: %s", err)
			}

			// this proposal is expired go for next proposal
			continue
		}

		result := proposal.DecideResult(&keyVault)
		switch result {
		case types.ProposalResult_PROPOSAL_RESULT_REJECTED:
			if err = k.finishPubkeysChangeProposal(
				ctx,
				proposal.Id,
				types.ProposalResult_PROPOSAL_RESULT_REJECTED,
				"rejected with more 'no' votes than 'yes' votes",
			); err != nil {
				return fmt.Errorf("error while setting the proposal as rejected: %s", err)
			}

			// this proposal is rejected go for next proposal
			continue
		case types.ProposalResult_PROPOSAL_RESULT_APPROVED:
			keyVault, found := k.GetKeyVault(ctx)
			if !found {
				fmt.Printf("there is no publick keys record")
			}

			keyVault.PublicKeys = proposal.Modifications.PublicKeys
			keyVault.SetLeader(proposal.Modifications.LeaderIndex)

			if err = k.finishPubkeysChangeProposal(
				ctx,
				proposal.Id,
				types.ProposalResult_PROPOSAL_RESULT_APPROVED,
				"",
			); err != nil {
				return fmt.Errorf("error while setting the proposal as approved: %s", err)
			}

			k.SetKeyVault(ctx, keyVault)
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
	proposal, found := k.GetPubkeysChangeProposal(
		ctx,
		types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
		proposalID,
	)
	if !found {
		return fmt.Errorf("proposal not found with id %d", proposalID)
	}

	proposal.Result = result
	proposal.ResultMeta = resultMetadata
	proposal.FinishTS = ctx.BlockTime().Unix()
	proposal.Status = types.ProposalStatus_PROPOSAL_STATUS_FINISHED

	k.RemoveProposal(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE, proposalID)
	k.SetPubkeysChangeProposal(ctx, proposal)

	return nil
}
