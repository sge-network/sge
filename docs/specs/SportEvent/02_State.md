# **State**

## **Params**

1. `event_min_bet_amount`: The minimum allowed bet amount that can be set in the whole system.
2. `event_min_bet_fee`: The minimum bet fee allowed across the system.
3. `event_max_sr_contribution`: The maximum allowed contribution by th sr module across the system.

```proto
// Params defines the parameters for the module.
// contains bet constraints associated to a sport-event.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // event_min_bet_amount default minimum bet amount for a sport-event.
  string event_min_bet_amount = 1 [
    (gogoproto.moretags) = "yaml:\"event_min_bet_amount\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // event_min_bet_fee default minimum bet fee for a sport-event.
  string event_min_bet_fee = 3 [
    (gogoproto.moretags) = "yaml:\"event_min_bet_fee\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // event_max_sr_contribution default max sr contribution for a sport-event.
  string event_max_sr_contribution = 4 [
    (gogoproto.moretags) = "yaml:\"event_max_sr_contribution\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

## **Bet Constraints**

```proto
// Bet constraints parent group for a sport-event
message EventBetConstraints {
  string min_amount = 2[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];


  cosmos.base.v1beta1.Coin bet_fee = 3
  [(gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin", (gogoproto.nullable) = false];
}
```

**EventBetConstraints**: Optional field which can put restrictions for bet acceptance criteria, this optional config help provide more
granular control of different categories of event. We can modify here following:

***MinBetAmount***: Minimum bet amount particular to this event.

***BetFee***: Fixed fee particular to the created event

---

## **Sport-Event**

```proto
// SportEvent the representation of the sport-event to be stored in
// the sport-event state.
message SportEvent {
  // uid is the universal unique identifier of the sport-event.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // start_ts is the start timestamp of the sport-event.
  uint64 start_ts = 2 [
    (gogoproto.customname) = "StartTS",
    (gogoproto.jsontag) = "start_ts",
    json_name = "start_ts"
  ];
  // end_ts is the end timestamp of the sport-event.
  uint64 end_ts = 3 [
    (gogoproto.customname) = "EndTS",
    (gogoproto.jsontag) = "end_ts",
    json_name = "end_ts"
  ];
  // odds is the list of odds of the sport-event.
  repeated Odds odds = 4;
  // winner_odds_uids is the list of winner odds universal unique identifiers.
  repeated string winner_odds_uids = 5 [
    (gogoproto.customname) = "WinnerOddsUIDs",
    (gogoproto.jsontag) = "winner_odds_uids",
    json_name = "winner_odds_uids"
  ];
  // status is the current status of the sport-event.
  SportEventStatus status = 6;
  // resolution_ts is the timestamp of the resolution of sport-event.
  uint64 resolution_ts = 7 [
    (gogoproto.customname) = "ResolutionTS",
    (gogoproto.jsontag) = "resolution_ts",
    json_name = "resolution_ts"
  ];
  // creator is the address of the creator of sport-event.
  string creator = 8;
  // bet_constraints holds the constraints of sport-event to accept bets.
  EventBetConstraints bet_constraints = 9;
  // meta contains human-readable metadata of the sport-event.
  string meta = 10;
  // sr_contribution_for_house is the amount of contibution of house in the sr
  string sr_contribution_for_house = 11 [
    (gogoproto.customname) = "SrContributionForHouse",
    (gogoproto.jsontag) = "sr_contribution_for_house",
    json_name = "sr_contribution_for_house",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // book_id is the id of book
  string book_id = 12 [
    (gogoproto.customname) = "BookID",
    (gogoproto.jsontag) = "book_id",
    json_name = "book_id"
  ];
}

// Bet constraints parent group for a sport-event
message EventBetConstraints {
  // min_amount is the minimum allowed bet amount for a sport-event.
  string min_amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // bet_fee is the fee that thebettor needs to pay to bet on the sport-event.
  string bet_fee = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

**UID**: universal unique sport-event ID.

**StartTS**: Timestamp when this event will start.

**EndTS**: Timestamp when this event will be over.

**Odds**: Array of all the associated values of type ***Odds*** with this event.

**WinnerOddsUids**: Array of all the Odds Uids which won this event, this would be a ***subset*** of the above provided OddsUidsArray.

**Status**: current state of the event.

**ResolutionTS**: Timestamp when the event came to a resolution i.e. we received a resolution request for an event.

**Creator**: Account responsible to create this event.

**BetConstraints**: Optional field which can put restrictions for bet acceptance criteria.

**Meta**: Human-Readable data of the event.

**SrContriButionForHouse** The amount percentage of contribution of sr in payouts

**BookID** The ID of the created order book

---

**type**: Enum

## **SportEventStatus**

```proto
// SportEventStatus is the sport-event status enumeration
enum SportEventStatus {
  // unspecified event
  SPORT_EVENT_STATUS_UNSPECIFIED = 0;
  // event is active
  SPORT_EVENT_STATUS_ACTIVE = 1;
  // event is inactive
  SPORT_EVENT_STATUS_INACTIVE = 2;
  // event is canceled
  SPORT_EVENT_STATUS_CANCELED = 3;
  // event is aborted
  SPORT_EVENT_STATUS_ABORTED = 4;
  // result of the event is declared
  SPORT_EVENT_STATUS_RESULT_DECLARED = 5;
}
```

---

## **Odds**

Is the type to represent odds item.

```proto
// Odds is a representation of and sport-event odds items.
message Odds {
  // uid is the universal unique identifier of the odds.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // meta contains any human-readable metadata of the odds.
  string meta = 2;
}
```

---

## **Statistics**

Keeps track of statistics of the sport-event module including the resolved unsettled sport-events.

```proto
// SportEventStats holds statistics of the sport-event
message SportEventStats {
  // resolved_unsettled is the list of universal unique identifiers
  // of resolved sport-events that have unsettled bets
  repeated string resolved_unsettled = 1;
}

```
