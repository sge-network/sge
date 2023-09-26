package types

import (
	context "context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
)

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
	definitions Definitions, ticket string,
) ([]Distribution, error) {
	var payload ApplyNoLossBetsRewardPayload
	if err := keepers.OVMKeeper.VerifyTicketUnmarshal(goCtx, ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(ErrInTicketVerification, "%s", err)
	}

	definition := definitions[0]

	if payload.Receiver.RecType != definition.RecType {
		return nil, sdkerrors.Wrapf(ErrAccReceiverTypeNotFound, "%s", payload.Receiver.RecType)
	}

	bettorAddr := payload.Receiver.Addr
	for _, betUID := range payload.BetUids {
		uID2ID, found := keepers.BetKeeper.GetBetID(ctx, betUID)
		if !found {
			return nil, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "bet id not found %s", betUID)
		}
		bet, found := keepers.BetKeeper.GetBet(ctx, bettorAddr, uID2ID.ID)
		if !found {
			return nil, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "bet not found %s", betUID)
		}
		if bet.Result != bettypes.Bet_RESULT_LOST {
			return nil, sdkerrors.Wrapf(ErrInvalidNoLossBetUID, "the bet result is not loss %s", betUID)
		}
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
