# **State**

## **Bet Constraints**

```
message EventBetConstraints {
  string max_bet_cap = 1[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  string min_amount = 2[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];


  cosmos.base.v1beta1.Coin bet_fee = 3
  [(gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin", (gogoproto.nullable) = false];

  string current_total_amount = 4[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  string max_loss = 5[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  string max_vig = 6[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable)   = false];

  string min_vig = 7[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable)   = false];

  map<string, uint64> current_bet_level_amount = 8;
}
```

**EventBetConstraints**: Optional field which can put restrictions for bet acceptance criteria, this optional config help provide more
granular control of different categories of event. We can modify here following:

***MaxBetCap***: Maximum overall capacity of the bet amount which can be accepted by this event during its lifetime.

***MinBetAmount***: Minimum bet amount particular to this event.

***BetFee***: Fixed fee particular to the created event

***CurrentTotalAmount***: This field is not meant to be exposed for alteration and would be maintained internally to set
current state of the particular sport event.

***MaxLoss***:This field will restrict us to accept more bets which supersedes the max loss for the event, this is an 
optional field and can be override while event creation/update

***Max/MinVig***:This field will be required primarily for bet module to calculate live vig to avoid placing any bad bet.
This is also an optional field and can be override while creation/update call

---

## **Sport Event**
```
message SportEvent {
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];
  uint64 start_ts = 2 [(gogoproto.customname) = "StartTS" ,(gogoproto.jsontag) = "start_ts", json_name = "start_ts"];
  uint64 end_ts = 3 [(gogoproto.customname) = "EndTS", (gogoproto.jsontag) = "end_ts", json_name = "end_ts"];
  repeated string odds_uids = 4 [(gogoproto.customname) = "OddsUIDs", (gogoproto.jsontag) = "odds_uids", json_name = "odds_uids"];
  map<string,bytes> winner_odds_uids = 5 [(gogoproto.customname) = "WinnerOddsUIDs", (gogoproto.jsontag) = "winner_odds_uids", json_name = "winner_odds_uids"];
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
	UID            string            `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid"`
	ResolutionTS   uint64            `protobuf:"varint,2,opt,name=resolution_ts,proto3" json:"resolution_ts"`
	WinnerOddsUIDs map[string][]byte `protobuf:"bytes,3,rep,name=winner_odds_uids,proto3" json:"winner_odds_uids" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Status         SportEventStatus  `protobuf:"varint,4,opt,name=status,proto3,enum=sgenetwork.sge.sportevent.SportEventStatus" json:"status,omitempty"`
}
```

---

## **UpdateEvent**

```
type UpdateEvent struct {
	Uid            string               `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	StartTS        uint64               `protobuf:"varint,2,opt,name=startTS,proto3" json:"startTS,omitempty"`
	EndTS          uint64               `protobuf:"varint,3,opt,name=endTS,proto3" json:"endTS,omitempty"`
	BetConstraints *EventBetConstraints `protobuf:"bytes,9,opt,name=betConstraints,proto3" json:"betConstraints,omitempty"`
	Active         bool                 `protobuf:"varint,10,opt,name=active,proto3" json:"active,omitempty"`
}
```

---

## **Standard Input message**
```
type MsgEvent struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Ticket  string `protobuf:"bytes,2,opt,name=ticket,proto3" json:"ticket,omitempty"`
}
```
- ticket is the standard dvm module entity

---

## **Standard Output message**
```
type MsgSportResponse struct {
	SuccessEvents []string       `protobuf:"bytes,1,rep,name=successEvents,proto3" json:"successEvents,omitempty"`
	FailedEvents  []*FailedEvent `protobuf:"bytes,2,rep,name=failedEvents,proto3" json:"failedEvents,omitempty"`
}
```
- success events will contain the ids of the successfully events generated
- failed events will contain the complete data of the failed input sent to the call
