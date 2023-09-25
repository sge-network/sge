package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
)

func NewDistribution(accAddr string, allocation Allocation) Distribution {
	return Distribution{
		AccAddr:    accAddr,
		Allocation: allocation,
	}
}

func NewAllocation(amount sdkmath.Int, receiverAccType ReceiverAccType, unlockTS uint64) Allocation {
	return Allocation{
		Amount:          amount,
		ReceiverAccType: receiverAccType,
		UnlockTS:        unlockTS,
	}
}

// ValidateBasic validates the basic properties of a reward definition.
func (d *Definition) ValidateBasic(blockTime uint64) error {
	if d.DstAccType != ReceiverAccType_RECEIVER_ACC_TYPE_SUB {
		if d.UnlockTS != 0 {
			return sdkerrors.Wrapf(ErrUnlockTSIsSubAccOnly, "%d", d.UnlockTS)
		}
	} else if d.UnlockTS <= blockTime {
		return sdkerrors.Wrapf(ErrUnlockTSDefBeforeBlockTime, "%d", d.UnlockTS)
	}
	return nil
}

// IRewardFactory defines the methods that should be implemented for all of reward types.
type IRewardFactory interface {
	VaidateDefinitions(campaign Campaign) error
	CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
		definitions []Definition, ticket string) ([]Distribution, error)
}

type RewardFactoryKeepers struct {
	OVMKeeper
	BetKeeper
}

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
		return sdkerrors.Wrapf(ErrInvalidReceiverType, "signup rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of signup reward.
func (sur SignUpReward) CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	definitions []Definition, ticket string,
) ([]Distribution, error) {
	var payload ApplySignupRewardPayload
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

// ReferralReward is the type for referral rewards calculations
type ReferralReward struct{}

// NewReferralReward create new object of referral reward calculator type.
func NewReferralReward() ReferralReward { return ReferralReward{} }

// VaidateDefinitions validates campaign definitions.
func (rfr ReferralReward) VaidateDefinitions(campaign Campaign) error {
	hasReferrer := false
	hasReferee := false
	for _, d := range campaign.RewardDefs {
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
	definitions []Definition, ticket string,
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
						d.DstAccType,
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

// NoLossBetsReward is the type for no loss bets rewards calculations
type NoLossBetsReward struct{}

// NewNoLossBetsReward create new object of no loss bets reward calculator type.
func NewNoLossBetsReward() NoLossBetsReward { return NoLossBetsReward{} }

// VaidateDefinitions validates campaign definitions.
func (afr NoLossBetsReward) VaidateDefinitions(campaign Campaign) error {
	if len(campaign.RewardDefs) != 1 {
		return sdkerrors.Wrapf(ErrWrongDefinitionsCount, "noloss bets rewards can only have single definition")
	}
	if campaign.RewardDefs[0].RecType != ReceiverType_RECEIVER_TYPE_SINGLE {
		return sdkerrors.Wrapf(ErrInvalidReceiverType, "noloss bets rewards can be defined for subaccount only")
	}
	return nil
}

// CalculateDistributions parses ticket payload and returns the distribution list of no loss bets reward.
func (afr NoLossBetsReward) CalculateDistributions(goCtx context.Context, ctx sdk.Context, keepers RewardFactoryKeepers,
	definitions []Definition, ticket string,
) ([]Distribution, error) {
	var payload ApplyNoLossBetsRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	bettorAddr := payload.Receiver.Addr
	for _, betUID := range payload.BetUids {
		betID, found := keepers.BetKeeper.GetBetID(ctx, betUID)
		if !found {
			return nil, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "bet id not found %s", betUID)
		}
		bet, found := keepers.BetKeeper.GetBet(ctx, bettorAddr, betID.ID)
		if !found {
			return nil, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "bet not founs %s", betUID)
		}
		if bet.Result != bettypes.Bet_RESULT_LOST {
			return nil, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "the bet result is not loss %s", betUID)
		}
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
