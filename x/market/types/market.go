package types

import (
	"github.com/mrz1836/go-sanitize"
)

func NewMarket(
	uid, creator string,
	startTS, endTS uint64,
	odds []*Odds,
	meta string,
	bookUID string,
	status MarketStatus,
) Market {
	return Market{
		UID:     uid,
		Creator: creator,
		StartTS: startTS,
		EndTS:   endTS,
		Odds:    odds,
		Meta:    sanitize.XSS(meta),
		BookUID: bookUID,
		Status:  status,
	}
}

// IsResolved returns true if the market is already resolved.
func (m *Market) IsResolved() bool {
	return m.Status == MarketStatus_MARKET_STATUS_RESULT_DECLARED ||
		m.Status == MarketStatus_MARKET_STATUS_CANCELED ||
		m.Status == MarketStatus_MARKET_STATUS_ABORTED
}

// IsUpdateAllowed returns true if updating the market is allowed.
func (m *Market) IsUpdateAllowed() bool {
	return m.isActiveOrInactive()
}

// IsResolveAllowed returns true if resolving the market is allowed.
func (m *Market) IsResolveAllowed() bool {
	return m.isActiveOrInactive()
}

func (m *Market) isActiveOrInactive() bool {
	return m.Status == MarketStatus_MARKET_STATUS_ACTIVE ||
		m.Status == MarketStatus_MARKET_STATUS_INACTIVE
}
