package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

// Validate validates fields of the given ticketData
func (payload *BetPlacementTicketPayload) Validate(creator string) error {
	if payload.SelectedOdds == nil {
		return ErrOddsDataNotFound
	}

	if !payload.KycData.Validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	if !utils.IsValidUID(payload.SelectedOdds.MarketUID) {
		return ErrInvalidMarketUID
	}

	if !utils.IsValidUID(payload.SelectedOdds.UID) {
		return ErrInvalidOddsUID
	}

	if len(strings.TrimSpace(payload.SelectedOdds.Value)) == 0 {
		return ErrEmptyOddsValue
	}

	if payload.SelectedOdds.MaxLossMultiplier.IsNil() ||
		payload.SelectedOdds.MaxLossMultiplier.LTE(sdk.ZeroDec()) {
		return ErrMaxLossMultiplierCanNotBeZero
	}

	if payload.SelectedOdds.MaxLossMultiplier.GT(sdk.OneDec()) {
		return ErrMaxLossMultiplierCanNotBeMoreThanOne
	}

	if payload.OddsType < OddsType_ODDS_TYPE_DECIMAL ||
		payload.OddsType > OddsType_ODDS_TYPE_MONEYLINE {
		return ErrInvalidOddsType
	}

	return nil
}
