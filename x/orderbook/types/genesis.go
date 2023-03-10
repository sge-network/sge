package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
)

// DefaultGenesis returns the default  genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                              DefaultParams(),
		BookList:                            []OrderBook{},
		BookParticipationList:               []BookParticipation{},
		BookExposureList:                    []BookOddsExposure{},
		ParticipationExposureList:           []ParticipationExposure{},
		ParticipationExposureByIndexList:    []ParticipationExposure{},
		HistoricalParticipationExposureList: []ParticipationExposure{},
		ParticipationBetPairExposureList:    []ParticipationBetPair{},
		PayoutLock:                          [][]byte{},
		Stats:                               OrderBookStats{ResolvedUnsettled: []string{}},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, p := range gs.BookParticipationList {
		_, err := sdk.AccAddressFromBech32(p.ParticipantAddress)
		if err != nil {
			return fmt.Errorf("invalid participant address %s", p.ParticipantAddress)
		}

		bookFound := false

		for _, b := range gs.BookList {
			if b.ID == p.BookUID {
				bookFound = true
			}
		}

		if !bookFound {
			return fmt.Errorf("book with id %s not found for participation %d", p.BookUID, p.Index)
		}
	}

	for _, b := range gs.BookList {
		oddsCount := 0
		for _, be := range gs.BookExposureList {
			exposureFound := false
			for _, pe := range gs.ParticipationExposureList {
				if pe.BookUID == b.ID && pe.OddsUID == be.OddsUID {
					exposureFound = true
				}
			}

			if !exposureFound {
				return fmt.Errorf("book with id %s not found for odds with uid %s", be.BookUID, be.OddsUID)
			}

			if be.BookUID == b.ID {
				oddsCount++
			}
		}

		exposureCount := cast.ToUint64(oddsCount)
		if exposureCount != b.OddsCount {
			return fmt.Errorf("book with id %s count does not match the odds exposure count %d", b.ID, exposureCount)
		}
	}

	for _, pei := range gs.ParticipationExposureByIndexList {
		exposureIndexFound := false
		for _, pe := range gs.ParticipationExposureList {
			if pei.BookUID == pe.BookUID &&
				pei.OddsUID == pe.OddsUID &&
				pei.ParticipationIndex == pe.ParticipationIndex {
				exposureIndexFound = true
			}
		}
		if !exposureIndexFound {
			return fmt.Errorf("participation index for the book %s, odds %s and index %d not found",
				pei.BookUID, pei.OddsUID, pei.ParticipationIndex)
		}
	}

	for _, pe := range gs.ParticipationExposureList {
		exposureIndexFound := false
		for _, pei := range gs.ParticipationExposureByIndexList {
			if pe.BookUID == pei.BookUID &&
				pe.OddsUID == pei.OddsUID &&
				pe.ParticipationIndex == pei.ParticipationIndex {
				exposureIndexFound = true
			}
		}
		if !exposureIndexFound {
			return fmt.Errorf("participation for the book %s, odds %s and index %d not found",
				pe.BookUID, pe.OddsUID, pe.ParticipationIndex)
		}
	}

	for _, hpei := range gs.HistoricalParticipationExposureList {
		exposureIndexFound := false
		for _, pe := range gs.ParticipationExposureList {
			if hpei.BookUID == pe.BookUID &&
				hpei.OddsUID == pe.OddsUID &&
				hpei.ParticipationIndex == pe.ParticipationIndex {
				exposureIndexFound = true
			}
		}
		if !exposureIndexFound {
			return fmt.Errorf("participation history for the book %s, odds %s and index %d not found",
				hpei.BookUID, hpei.OddsUID, hpei.ParticipationIndex)
		}
	}

	for _, bp := range gs.ParticipationBetPairExposureList {
		bookFound := false
		for _, b := range gs.BookList {
			if bp.BookUID == b.ID {
				bookFound = true
			}
		}
		if !bookFound {
			return fmt.Errorf("participation bet pair not found for the book %s", bp.BookUID)
		}
	}

	return gs.Params.Validate()
}
