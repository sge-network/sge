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
// the internal keeper representation of sport-event
message SportEvent {
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];
  uint64 start_ts = 2 [(gogoproto.customname) = "StartTS" ,(gogoproto.jsontag) = "start_ts", json_name = "start_ts"];
  uint64 end_ts = 3 [(gogoproto.customname) = "EndTS", (gogoproto.jsontag) = "end_ts", json_name = "end_ts"];
  repeated Odds odds = 4;
  repeated string winner_odds_uids = 5 [(gogoproto.customname) = "WinnerOddsUIDs", (gogoproto.jsontag) = "winner_odds_uids", json_name = "winner_odds_uids"];
  SportEventStatus status = 6;
  uint64 resolution_ts = 7 [(gogoproto.customname) = "ResolutionTS", (gogoproto.jsontag) = "resolution_ts", json_name = "resolution_ts"];
  string creator = 8;
  EventBetConstraints betConstraints = 9;
  bool active = 10;
  string meta = 11;
}
```

**Uid**: unique event Id.

**StartTS**:Time when this event will start

**EndTS**: Time when this event will be over

**Odds**: Array of all the associated values of type ***Odds*** with this event.

**WinnerOddsUids**: Array of all the Odds Uids which won this event, this would be a ***subset*** of the above provided OddsUidsArray

**Status**: current state of the event.

**ResolutionTS**: Time when the event came to a resolution i.e. we received a resolution request for an event

**Creator**: Account responsible to create this event

**EventBetConstraints**: Optional field which can put restrictions for bet acceptance criteria.

**Active**: This is a toggle field which marks if the event is active(open to take or accepts bets)

**Meta**: Human-Readable data of the event.

---

## **SportEventStatus**

```go
(
 SportEventStatus_STATUS_PENDING         SportEventStatus = 0
 SportEventStatus_STATUS_INVALID         SportEventStatus = 1
 SportEventStatus_STATUS_CANCELLED       SportEventStatus = 2
 SportEventStatus_STATUS_ABORTED         SportEventStatus = 3
 SportEventStatus_STATUS_RESULT_DECLARED SportEventStatus = 4
)
```

**type**: Enum

**total_types**: 5

---

## **ResolutionEvent**

```go
type ResolutionEvent struct {
 UID            string            `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid"`
 ResolutionTS   uint64            `protobuf:"varint,2,opt,name=resolution_ts,proto3" json:"resolution_ts"`
 WinnerOddsUIDs []string `protobuf:"bytes,3,rep,name=winner_odds_uids,proto3" json:"winner_odds_uids" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
 Status         SportEventStatus  `protobuf:"varint,4,opt,name=status,proto3,enum=sgenetwork.sge.sportevent.SportEventStatus" json:"status,omitempty"`
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
