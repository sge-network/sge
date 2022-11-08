# **State**

## **Bet Constraints**

```
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

## **Sport Event**
```
message SportEvent {
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];
  uint64 start_ts = 2 [(gogoproto.customname) = "StartTS" ,(gogoproto.jsontag) = "start_ts", json_name = "start_ts"];
  uint64 end_ts = 3 [(gogoproto.customname) = "EndTS", (gogoproto.jsontag) = "end_ts", json_name = "end_ts"];
  repeated string odds_uids = 4 [(gogoproto.customname) = "OddsUIDs", (gogoproto.jsontag) = "odds_uids", json_name = "odds_uids"];
  repeated string winner_odds_uids = 5 [(gogoproto.customname) = "WinnerOddsUIDs", (gogoproto.jsontag) = "winner_odds_uids", json_name = "winner_odds_uids"];
  SportEventStatus status = 6;
  uint64 resolution_ts = 7 [(gogoproto.customname) = "ResolutionTS", (gogoproto.jsontag) = "resolution_ts", json_name = "resolution_ts"];
  string creator = 8;
  EventBetConstraints betConstraints = 9;
  bool active = 10;
}
```
**Uid**: unique event Id.

**StartTs**:Time when this event will start

**EndTs**: Time when this event will be over

**OddsUids**: Array of all the associated **_odds_** with this event ex: (uid_win_event, uid_draw_event, uid_loose_event)

**WinnerOddsUids**: Array of all the Odds Uids which won this event, this would be a **_subset_** of the above provided OddsUidsArray

**Status**: current state of the event.

**ResolutionTs**: Time when the event came to a resolution i.e. we received a resolution request for an event

**Creator**: Account responsible to create this event

**Active**: This is a toggle field which marks if the event is active(open to take or accepts bets)

---

## **SportEventStatus**

```
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

```
type ResolutionEvent struct {
	string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];
    uint64 resolution_ts = 2 [(gogoproto.customname) = "ResolutionTS", (gogoproto.jsontag) = "resolution_ts", json_name = "resolution_ts"];
    repeated string winner_odds_uids = 3 [(gogoproto.customname) = "WinnerOddsUIDs", (gogoproto.jsontag) = "winner_odds_uids", json_name = "winner_odds_uids"];
    SportEventStatus status = 4;
}
```
