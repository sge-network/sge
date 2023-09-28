package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SignUpReward is the type for signup rewards calculations
type SignUpReward struct{}

// NewSignUpReward create new object of signup reward calculator type.
func NewSignUpReward() SignUpReward { return SignUpReward{} }

// VaidateDefinitions validates campaign definitions.
func (sur SignUpReward) VaidateDefinitions(campaign Campaign) error {
	if len(campaign.RewardDefs) != 1 {
		return sdkerrors.Wrapf(ErrWrongDefinitionsCount, "signup rewards can only have single definition")
	}
	if campaign.RewardDefs[0].RecType != ReceiverType_RECEIVER_TYPE_SINGLE {
		return sdkerrors.Wrapf(ErrInvalidReceiverType, "signup rewards can be defined for single receiver only")
	}
	if campaign.RewardDefs[0].DstAccType != ReceiverAccType_RECEIVER_ACC_TYPE_SUB {
		return sdkerrors.Wrapf(ErrInvalidReceiverType, "signup rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpReward) CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	definitions Definitions, ticket string,
) ([]Distribution, error) {
	var payload ApplySignupRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	definition := definitions[0]

	if payload.Receiver.RecType != definition.RecType {
		return nil, sdkerrors.Wrapf(ErrAccReceiverTypeNotFound, "%s", payload.Receiver.RecType)
	}

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
