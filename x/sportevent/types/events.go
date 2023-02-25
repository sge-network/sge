package types

const (
	AttributeKeyEventsCreator = "events_creator"

	AttributeKeySportEventsSuccessUID = "sport_events_success_uid"
	AttributeKeyOrderBookUID          = "sport_events_book_uid"
	AttributeKeySportEventsFailedUID  = "sport_events_failed_uid"

	// AttributeValueCategory is the event attribute for category as module name
	AttributeValueCategory = ModuleName
)

const (
	TypeMsgCreateSportEvents  = "create_sport_events"
	TypeMsgUpdateSportEvents  = "update_sport_events"
	TypeMsgResolveSportEvents = "resolve_sport_events"
)
