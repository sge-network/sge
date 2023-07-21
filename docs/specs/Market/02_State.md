# **State**

## **Params**

The market module does not have any parameters.

---

## **Market**

```proto
// Market is the representation of the market to be stored in
// the market state.
message Market {
  // uid is the universal unique identifier of the market.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // start_ts is the start timestamp of the market.
  uint64 start_ts = 2 [
    (gogoproto.customname) = "StartTS",
    (gogoproto.jsontag) = "start_ts",
    json_name = "start_ts"
  ];
  // end_ts is the end timestamp of the market.
  uint64 end_ts = 3 [
    (gogoproto.customname) = "EndTS",
    (gogoproto.jsontag) = "end_ts",
    json_name = "end_ts"
  ];
  // odds is the list of odds of the market.
  repeated Odds odds = 4;
  // winner_odds_uids is the list of winner odds universal unique identifiers.
  repeated string winner_odds_uids = 5 [
    (gogoproto.customname) = "WinnerOddsUIDs",
    (gogoproto.jsontag) = "winner_odds_uids",
    json_name = "winner_odds_uids"
  ];
  // status is the current status of the market.
  MarketStatus status = 6;
  // resolution_ts is the timestamp of the resolution of market.
  uint64 resolution_ts = 7 [
    (gogoproto.customname) = "ResolutionTS",
    (gogoproto.jsontag) = "resolution_ts",
    json_name = "resolution_ts"
  ];
  // creator is the address of the creator of market.
  string creator = 8;
  // meta contains human-readable metadata of the market.
  string meta = 9;
  // book_uid is the unique identifier corresponding to the book
  string book_uid = 10 [
    (gogoproto.customname) = "BookUID",
    (gogoproto.jsontag) = "book_uid",
    json_name = "book_uid"
  ];
}
```

**UID**: universal unique market ID.

**StartTS**: Timestamp when this market will start.

**EndTS**: Timestamp when this market will be over.

**Odds**: Array of all the associated values of type ***Odds*** with this market.

**WinnerOddsUids**: Array of all the Odds Uids which won this market, this would be a ***subset*** of the above provided OddsUidsArray.

**Status**: current state of the market.

**ResolutionTS**: Timestamp when the market came to a resolution i.e. we received a resolution request for an market.

**Creator**: Account responsible to create this market.

**Meta**: Human-Readable data of the market.

**BookID** The ID of the created order book

---

**type**: Enum

## **MarketStatus**

```proto
// MarketStatus is the market status enumeration
enum MarketStatus {
  // unspecified market
  MARKET_STATUS_UNSPECIFIED = 0;
  // market is active
  MARKET_STATUS_ACTIVE = 1;
  // market is inactive
  MARKET_STATUS_INACTIVE = 2;
  // market is canceled
  MARKET_STATUS_CANCELED = 3;
  // market is aborted
  MARKET_STATUS_ABORTED = 4;
  // result of the market is declared
  MARKET_STATUS_RESULT_DECLARED = 5;
}
```

---

## **Odds**

Is the type to represent odds item.

```proto
// Odds is a representation of market odds.
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

Keeps track of statistics of the market module including the resolved unsettled markets.

```proto
// MarketStats holds statistics of the market
message MarketStats {
  // resolved_unsettled is the list of universal unique identifiers
  // of resolved markets that have unsettled bets
  repeated string resolved_unsettled = 1;
}

```
