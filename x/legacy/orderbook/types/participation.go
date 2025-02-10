package types

import (
	yaml "gopkg.in/yaml.v2"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	housetypes "github.com/sge-network/sge/x/legacy/house/types"
)

// NewOrderBookParticipation creates a new book participation object
//
//nolint:interfacer
func NewOrderBookParticipation(
	index uint64,
	orderBookUID string,
	participantAddress string,
	exposuresNotFilled uint64,
	liquidity, fee, currentRoundLiquidity, totalBetAmount, currentRoundTotalBetAmount, maxLoss, currentRoundMaxLoss sdkmath.Int,
	currentRoundMaxLossOddsUID string,
	actualProfit sdkmath.Int,
) OrderBookParticipation {
	return OrderBookParticipation{
		Index:                      index,
		OrderBookUID:               orderBookUID,
		ParticipantAddress:         participantAddress,
		Liquidity:                  liquidity,
		Fee:                        fee,
		CurrentRoundLiquidity:      currentRoundLiquidity,
		ExposuresNotFilled:         exposuresNotFilled,
		TotalBetAmount:             totalBetAmount,
		CurrentRoundTotalBetAmount: currentRoundTotalBetAmount,
		MaxLoss:                    maxLoss,
		CurrentRoundMaxLoss:        currentRoundMaxLoss,
		CurrentRoundMaxLossOddsUID: currentRoundMaxLossOddsUID,
		ActualProfit:               actualProfit,
	}
}

// String returns a human readable string representation of a BookParticipation.
func (p OrderBookParticipation) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// CalculateMaxLoss calculates the maxixmum amount of the tokens expected to be the
// loss of the participation according to the bet amount
func (p OrderBookParticipation) CalculateMaxLoss(betAmount sdkmath.Int) sdkmath.Int {
	return p.CurrentRoundMaxLoss.Sub(betAmount)
}

// ValidateWithdraw determines if the participation is allowed
// to be withdrawn or not.
func (p *OrderBookParticipation) ValidateWithdraw(
	depositorAddress string,
	participationIndex uint64,
) error {
	if p.IsSettled {
		return sdkerrors.Wrapf(
			ErrBookParticipationAlreadySettled,
			"%s, %d",
			p.OrderBookUID,
			participationIndex,
		)
	}

	if p.ParticipantAddress != depositorAddress {
		return sdkerrors.Wrapf(ErrMismatchInDepositorAddress, "%s", p.ParticipantAddress)
	}

	return nil
}

// maxWithdrawalAmount returns the max withdrawal amount of a participation.
func (p *OrderBookParticipation) maxWithdrawalAmount() sdkmath.Int {
	if p.CurrentRoundMaxLoss.LT(sdkmath.ZeroInt()) {
		return p.CurrentRoundLiquidity
	}
	return p.CurrentRoundLiquidity.Sub(p.CurrentRoundMaxLoss)
}

// IsLiquidityInCurrentRound determines if the participation has liquidity in current round.
func (p *OrderBookParticipation) IsLiquidityInCurrentRound() bool {
	return p.CurrentRoundLiquidity.GT(sdkmath.ZeroInt())
}

// WithdrawableAmount returns the withdrawal amount according to the withdrawal mode and max withdrawable amount.
func (p *OrderBookParticipation) WithdrawableAmount(
	mode housetypes.WithdrawalMode,
	amount sdkmath.Int,
) (sdkmath.Int, error) {
	// Calculate max amount that can be transferred
	maxTransferableAmount := p.maxWithdrawalAmount()

	var withdrawalAmt sdkmath.Int
	switch mode {
	case housetypes.WithdrawalMode_WITHDRAWAL_MODE_FULL:
		if maxTransferableAmount.LTE(sdkmath.ZeroInt()) {
			return sdkmath.Int{}, sdkerrors.Wrapf(
				ErrMaxWithdrawableAmountIsZero,
				"%d, %d",
				p.CurrentRoundLiquidity,
				p.CurrentRoundMaxLoss,
			)
		}
		withdrawalAmt = maxTransferableAmount
	case housetypes.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL:
		if maxTransferableAmount.LT(amount) {
			return sdkmath.Int{}, sdkerrors.Wrapf(
				ErrWithdrawalAmountIsTooLarge,
				": got %s, max %s",
				amount,
				maxTransferableAmount,
			)
		}
		withdrawalAmt = amount
	default:
		return sdkmath.Int{}, sdkerrors.Wrapf(housetypes.ErrInvalidMode, "%s", mode.String())
	}

	return withdrawalAmt, nil
}

// SetLiquidityAfterWithdrawal sets the liquidity props after withdrawal.
func (p *OrderBookParticipation) SetLiquidityAfterWithdrawal(withdrawalAmt sdkmath.Int) {
	p.CurrentRoundLiquidity = p.CurrentRoundLiquidity.Sub(withdrawalAmt)
	p.Liquidity = p.Liquidity.Sub(withdrawalAmt)
}

// NotParticipatedInBetFulfillment determines if the participation has
// participated in the bet fulfillment.
func (p *OrderBookParticipation) NotParticipatedInBetFulfillment() bool {
	return p.TotalBetAmount.Equal(sdkmath.ZeroInt())
}

// IsEligibleForNextRound determines if the participation has enough
// liquidity to be used in the next round or not
func (p *OrderBookParticipation) IsEligibleForNextRound() bool {
	return p.CurrentRoundLiquidity.GT(sdkmath.ZeroInt())
}

// IsEligibleForNextRound determines if the participation has enough
// liquidity to be used in the next round or not
func (p *OrderBookParticipation) IsEligibleForNextRoundPreLiquidityReduction() bool {
	maxLoss := sdkmath.MaxInt(sdkmath.ZeroInt(), p.CurrentRoundMaxLoss)
	return p.CurrentRoundLiquidity.Sub(maxLoss).GT(sdkmath.ZeroInt())
}

// TrimCurrentRoundLiquidity subtracts the max loss from the current round liquidity.
func (p *OrderBookParticipation) TrimCurrentRoundLiquidity() {
	maxLoss := sdkmath.MaxInt(sdkmath.ZeroInt(), p.CurrentRoundMaxLoss)
	p.CurrentRoundLiquidity = p.CurrentRoundLiquidity.Sub(maxLoss)
}

// ResetForNextRound resets the exposures, max loss and current round amount
// and make the participation ready for the next round
func (p *OrderBookParticipation) ResetForNextRound(notFilledExposures uint64) {
	// prepare participation for the next round
	p.ExposuresNotFilled = notFilledExposures
	p.MaxLoss = p.MaxLoss.Add(p.CurrentRoundMaxLoss)
	p.CurrentRoundTotalBetAmount = sdkmath.ZeroInt()
	p.CurrentRoundMaxLoss = sdkmath.ZeroInt()
}

// SetCurrentRound sets the current round total bet amount and max loss.
func (p *OrderBookParticipation) SetCurrentRound(
	pe *ParticipationExposure,
	oddsUID string,
	betAmount sdkmath.Int,
) {
	p.TotalBetAmount = p.TotalBetAmount.Add(betAmount)
	p.CurrentRoundTotalBetAmount = p.CurrentRoundTotalBetAmount.Add(betAmount)
	p.setMaxLoss(pe, oddsUID, betAmount)
}

// setMaxLoss sets the max loss of the current round.
func (p *OrderBookParticipation) setMaxLoss(
	pe *ParticipationExposure,
	oddsUID string,
	betAmount sdkmath.Int,
) {
	// max loss is the maximum amount that an exposure may lose.
	maxLoss := pe.CalculateMaxLoss(p.CurrentRoundTotalBetAmount)
	switch {
	case p.CurrentRoundMaxLoss.IsNil():
		p.CurrentRoundMaxLoss = maxLoss
		p.CurrentRoundMaxLossOddsUID = oddsUID
	case p.CurrentRoundMaxLossOddsUID == oddsUID:
		p.CurrentRoundMaxLoss = maxLoss
	default:
		originalMaxLoss := p.CalculateMaxLoss(betAmount)
		if maxLoss.GT(originalMaxLoss) {
			p.CurrentRoundMaxLoss = maxLoss
			p.CurrentRoundMaxLossOddsUID = oddsUID
		} else {
			p.CurrentRoundMaxLoss = originalMaxLoss
		}
	}
}
