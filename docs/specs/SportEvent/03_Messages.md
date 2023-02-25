# **Messages**

The Sportevent module exposes the following services:

- AddEvent
- ResolveEvent
- UpdateEvent

```proto
// Msg defines the Msg service.
service Msg {
    rpc AddSportEvent(MsgAddSportEvent) returns (SportEventResponse);
    rpc ResolveSportEvent(MsgResolveSportEvent) returns (SportEventResponse);
    rpc UpdateSportEvent(MsgUpdateSportEvent) returns (SportEventResponse);
}
```

---

## **MsgAddSportEvent**

This message is used to add one or more new sportevent to the chain

```proto
message MsgAddSportEvent {
  string creator = 1;
  string ticket = 2;
}
```

### Add Sport Event Ticked Payload

The ticket data for sport-event addition is as follows:

```proto
// SportEventAddTicketPayload indicates data of add sport-event ticket
message SportEventAddTicketPayload {
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

  // status is the current status of the sport-event.
  SportEventStatus status = 5;

  // creator is the address of the creator of sport-event.
  string creator = 6;

  // min_bet_amount is the minimum allowed bet amount for a sport-event.
  string min_bet_amount = 7 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // bet_fee is the fee that thebettor needs to pay to bet on the sport-event.
  string bet_fee = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // meta contains human-readable metadata of the sport-event.
  string meta = 9;
  // sr_contribution_for_house is the portion of payout that should be paid by
  // sr
  string sr_contribution_for_house = 10 [
    (gogoproto.customname) = "SrContributionForHouse",
    (gogoproto.jsontag) = "sr_contribution_for_house",
    json_name = "sr_contribution_for_house",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

#### **Sample addition ticket**

```json
{
    "uid": "5531c60f-2025-48ce-ae79-1dc110f16000",
    "start_ts": 1668480139,
    "end_ts": 1883781609,
    "odds": [
        {
            "uid": "9991c60f-2025-48ce-ae79-1dc110f16990",
            "meta": "x is winner"
        },
        {
            "uid": "9991c60f-2025-48ce-ae79-1dc110f16991",
            "meta": "y is winner"
        },
        {
            "uid": "9991c60f-2025-48ce-ae79-1dc110f16992",
            "meta": "draw"
        }
    ],
    "status": 1,
    "min_bet_amount": "1000000",
    "bet_fee": "10",
    "meta": "Soccer: England vs USA",
    "sr_contribution_for_house": "5",
    "iat": 1665140310,
    "exp": 1757788212
}
```

---

## **MsgResolveSportEvent**

This message is used to resolve one or more already existent events on the chain

```proto
message MsgResolveSportEvent {
  string creator = 1;
  string ticket = 2;
}
```

### Resolve Sport Event Ticked Payload

```proto
// SportEventResolutionTicketPayload indicates data of the
// resolution of the sport-event ticket.
message SportEventResolutionTicketPayload {
  // uid is the universal unique identifier of sport-event.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // resolution_ts is the resolution timestamp of the sport-event.
  uint64 resolution_ts = 2 [
    (gogoproto.customname) = "ResolutionTS",
    (gogoproto.jsontag) = "resolution_ts",
    json_name = "resolution_ts"
  ];

  // winner_odds_uids is the universal unique identifier list of the winner
  // odds.
  repeated string winner_odds_uids = 3 [
    (gogoproto.customname) = "WinnerOddsUIDs",
    (gogoproto.jsontag) = "winner_odds_uids",
    json_name = "winner_odds_uids"
  ];

  // status is the status of the resolution.
  SportEventStatus status = 4;
}
```

#### **Sample resolve ticket**

```json
{
    "uid": "5531c60f-2025-48ce-ae79-1dc110f16000",
    "resolution_ts": 1668480139,
    "winner_odds_uids": [
      "9991c60f-2025-48ce-ae79-1dc110f16990"
    ],
    "status": 1,
    "iat": 1665140310,
    "exp": 1757788212
}
```

---

## **MsgUpdateSportEvent**

This message is used to update one or more already existent events on the chain

```proto
message MsgUpdateSportEvent {
  string creator = 1;
  string ticket = 2;
}
```

### Update Sport Event Ticked Payload

```proto
// SportEventUpdateTicketPayload indicates data of update sport-event ticket
message SportEventUpdateTicketPayload {
  // uid is the uuid of the sport-event
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // start_ts is the start timestamp of the sport-event
  uint64 start_ts = 2 [
    (gogoproto.customname) = "StartTS",
    (gogoproto.jsontag) = "start_ts",
    json_name = "start_ts"
  ];
  // end_ts is the end timestamp of the sport-event
  uint64 end_ts = 3 [
    (gogoproto.customname) = "EndTS",
    (gogoproto.jsontag) = "end_ts",
    json_name = "end_ts"
  ];
  // min_bet_amount is the minimum allowed bet amount for a sport-event.
  string min_bet_amount = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // bet_fee is the fee that thebettor needs to pay to bet on the sport-event.
  string bet_fee = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // status is the status of the resolution.
  SportEventStatus status = 6;
}
```

#### **Sample update ticket**

```json
{
    "uid": "5531c60f-2025-48ce-ae79-1dc110f16000",
    "start_ts": 1668480139,
    "end_ts": 1883781609,
    "status": 1,
    "min_bet_amount": "1000000",
    "bet_fee": "10",
    "iat": 1665140310,
    "exp": 1757788212
}
```

---

## **SportEventResponse**

This is the common response to all the messages

```proto
// MsgAddSportEventResponse response for adding sport-event.
message MsgAddSportEventResponse {
  // error contains an error if adding a sport-event faces any issues.
  string error = 1 [ (gogoproto.nullable) = true ];
  // data is the data of sport-event.
  SportEvent data = 2 [ (gogoproto.nullable) = true ];
}
```
