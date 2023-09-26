package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AffiliationReward is the type for affiliation rewards calculations
type AffiliationReward struct{}

// NewAffiliationReward create new object of affiliation reward calculator type.
func NewAffiliationReward() AffiliationReward { return AffiliationReward{} }

// VaidateDefinitions validates campaign definitions.
func (afr AffiliationReward) VaidateDefinitions(campaign Campaign) error {
	if len(campaign.RewardDefs) != 1 {
		return sdkerrors.Wrapf(ErrWrongDefinitionsCount, "affiliation rewards can only have single definition")
	}
	if campaign.RewardDefs[0].RecType != ReceiverType_RECEIVER_TYPE_SINGLE {
		return sdkerrors.Wrapf(ErrInvalidReceiverType, "affiliation rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of affiliation reward.
func (afr AffiliationReward) CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	definitions []Definition, ticket string,
) ([]Distribution, error) {
	var payload ApplyAffiliationRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	definition := definitions[0]

	return []Distribution{
		NewDistribution(
			payload.Receiver.Addr,
			NewAllocation(
				definition.Amount,
				definition.DstAccType,
				definition.UnlockTS,
			),
		),
	}, nil
}
