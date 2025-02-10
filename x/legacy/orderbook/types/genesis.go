package types

import (
	"fmt"

	"github.com/spf13/cast"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                              DefaultParams(),
		OrderBookList:                       []OrderBook{},
		OrderBookParticipationList:          []OrderBookParticipation{},
		OrderBookExposureList:               []OrderBookOddsExposure{},
		ParticipationExposureList:           []ParticipationExposure{},
		ParticipationExposureByIndexList:    []ParticipationExposure{},
		HistoricalParticipationExposureList: []ParticipationExposure{},
		ParticipationBetPairExposureList:    []ParticipationBetPair{},
		Stats:                               OrderBookStats{ResolvedUnsettled: []string{}},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, p := range gs.OrderBookParticipationList {
		_, err := sdk.AccAddressFromBech32(p.ParticipantAddress)
		if err != nil {
			return fmt.Errorf("invalid participant address %s", p.ParticipantAddress)
		}

		bookFound := false

		for _, b := range gs.OrderBookList {
			if b.UID == p.OrderBookUID {
				bookFound = true
			}
		}

		if !bookFound {
			return fmt.Errorf("book with id %s not found for participation %d", p.OrderBookUID, p.Index)
		}
	}

	for _, b := range gs.OrderBookList {
		oddsCount := 0
		for _, be := range gs.OrderBookExposureList {
			exposureFound := false
			for _, pe := range gs.ParticipationExposureList {
				if pe.OrderBookUID == b.UID && pe.OddsUID == be.OddsUID {
					exposureFound = true
				}
			}

			if !exposureFound {
				return fmt.Errorf(
					"book with id %s not found for odds with uid %s",
					be.OrderBookUID,
					be.OddsUID,
				)
			}

			if be.OrderBookUID == b.UID {
				oddsCount++
			}
		}

		exposureCount := cast.ToUint64(oddsCount)
		if exposureCount != b.OddsCount {
			return fmt.Errorf(
				"book with id %s count does not match the odds exposure count %d",
				b.UID,
				exposureCount,
			)
		}
	}

	for _, pei := range gs.ParticipationExposureByIndexList {
		exposureIndexFound := false
		for _, pe := range gs.ParticipationExposureList {
			if pei.OrderBookUID == pe.OrderBookUID &&
				pei.OddsUID == pe.OddsUID &&
				pei.ParticipationIndex == pe.ParticipationIndex {
				exposureIndexFound = true
			}
		}
		if !exposureIndexFound {
			return fmt.Errorf("participation index for the book %s, odds %s and index %d not found",
				pei.OrderBookUID, pei.OddsUID, pei.ParticipationIndex)
		}
	}

	for _, pe := range gs.ParticipationExposureList {
		exposureIndexFound := false
		for _, pei := range gs.ParticipationExposureByIndexList {
			if pe.OrderBookUID == pei.OrderBookUID &&
				pe.OddsUID == pei.OddsUID &&
				pe.ParticipationIndex == pei.ParticipationIndex {
				exposureIndexFound = true
			}
		}
		if !exposureIndexFound {
			return fmt.Errorf("participation for the book %s, odds %s and index %d not found",
				pe.OrderBookUID, pe.OddsUID, pe.ParticipationIndex)
		}
	}

	for _, hpei := range gs.HistoricalParticipationExposureList {
		exposureIndexFound := false
		for _, pe := range gs.ParticipationExposureList {
			if hpei.OrderBookUID == pe.OrderBookUID &&
				hpei.OddsUID == pe.OddsUID &&
				hpei.ParticipationIndex == pe.ParticipationIndex {
				exposureIndexFound = true
			}
		}
		if !exposureIndexFound {
			return fmt.Errorf("participation history for the book %s, odds %s and index %d not found",
				hpei.OrderBookUID, hpei.OddsUID, hpei.ParticipationIndex)
		}
	}

	for _, bp := range gs.ParticipationBetPairExposureList {
		bookFound := false
		for _, b := range gs.OrderBookList {
			if bp.OrderBookUID == b.UID {
				bookFound = true
			}
		}
		if !bookFound {
			return fmt.Errorf("participation bet pair not found for the book %s", bp.OrderBookUID)
		}
	}

	return gs.Params.Validate()
}
