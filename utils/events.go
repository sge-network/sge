package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewEventEmitter(ctx *sdk.Context, category string) EventEmitter {
	return EventEmitter{ctx: ctx, category: category}
}

type EventEmitter struct {
	ctx      *sdk.Context
	category string
	Events   []sdk.Event
}

func (em *EventEmitter) AddMsg(msgType, sender string, attrs ...sdk.Attribute) {
	em.addEvent(sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, em.category),
		sdk.NewAttribute(sdk.AttributeKeyAction, msgType),
		sdk.NewAttribute(sdk.AttributeKeySender, sender),
	)
	em.addEvent(msgType, attrs...)
}

func (em *EventEmitter) addEvent(ty string, attrs ...sdk.Attribute) {
	em.Events = append(em.Events, sdk.NewEvent(ty, attrs...))
}

func (em EventEmitter) Emit() {
	em.ctx.EventManager().EmitEvents(em.Events)
}
