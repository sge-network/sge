package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
)

func NewOrderBookEvent() OrderBookEvent {
	return OrderBookEvent{
		ParticipationExposure: []*ParticipationExposure{},
		OrderBookOddsExposure: []*OrderBookOddsExposure{},
	}
}

func (obe *OrderBookEvent) AddParticipationExposure(pe ParticipationExposure) {
	obe.ParticipationExposure = append(obe.ParticipationExposure, &pe)
}

func (obe *OrderBookEvent) AddOrderBookOddsExposure(boe OrderBookOddsExposure) {
	obe.OrderBookOddsExposure = append(obe.OrderBookOddsExposure, &boe)
}

func (obe *OrderBookEvent) Emit(ctx sdk.Context) {
	emitter := utils.NewEventEmitter(&ctx, AttributeKeyCategoryOrderBookEvent)
	emitter.AddEvent(EventTypeOrderbook,
		sdk.NewAttribute(AttributeOrderBookEvent, obe.String()),
	)
	emitter.Emit()
}
