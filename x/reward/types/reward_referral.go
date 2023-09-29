package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReferralReward is the type for referral rewards calculations
type ReferralReward struct{}

// NewReferralReward create new object of referral reward calculator type.
func NewReferralReward() ReferralReward { return ReferralReward{} }

// VaidateDefinitions validates campaign definitions.
func (rfr ReferralReward) VaidateDefinitions(campaign Campaign) error {
	hasReferrer := false
	hasReferee := false
	for _, d := range campaign.RewardDefs {
		if d.RecAccType != ReceiverAccType_RECEIVER_ACC_TYPE_SUB {
			return sdkerrors.Wrapf(ErrInvalidReceiverType, "referral rewards can be defined for subaccount only")
		}
		switch d.RecType {
		case ReceiverType_RECEIVER_TYPE_REFEREE:
			hasReferee = true
		case ReceiverType_RECEIVER_TYPE_REFERRER:
			hasReferrer = true
		default:
			return sdkerrors.Wrapf(ErrInvalidReceiverType, "%s", d.RecType)
		}
	}

	if !hasReferee || !hasReferrer {
		return sdkerrors.Wrapf(ErrMissingDefinition, "referral rewards should have the referrer and the referee")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of referral reward.
func (rfr ReferralReward) CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	definitions Definitions, ticket string,
) ([]Distribution, error) {
	var payload ApplyRerferralRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	distributions := []Distribution{}
	for _, d := range definitions {
		found := false
		for _, r := range payload.Receivers {
			if d.RecType == r.RecType {
				found = true
				distributions = append(distributions, NewDistribution(
					r.Addr,
					NewAllocation(
						d.Amount,
						d.RecAccType,
						d.UnlockTS,
					),
				))
			}
		}
		if !found {
			return nil, sdkerrors.Wrapf(ErrAccReceiverTypeNotFound, "%s", d.RecType)
		}
	}

	return distributions, nil
}
