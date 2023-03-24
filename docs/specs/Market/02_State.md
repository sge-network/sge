# **State**

## **Params**

1. `min_bet_amount`: The minimum allowed bet amount that can be set in the whole system.
2. `min_bet_fee`: The minimum bet fee allowed across the system.
3. `max_sr_contribution`: The maximum allowed contribution by th sr module across the system.

```proto
// Params defines the parameters for the module.
// contains bet constraints associated to a market.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // min_bet_amount default minimum bet amount for a market.
  string min_bet_amount = 1 [
    (gogoproto.moretags) = "yaml:\"min_bet_amount\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // min_bet_fee default minimum bet fee for a market.
  string min_bet_fee = 3 [
    (gogoproto.moretags) = "yaml:\"min_bet_fee\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // max_sr_contribution default max sr contribution for a market.
  string max_sr_contribution = 4 [
    (gogoproto.moretags) = "yaml:\"max_sr_contribution\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

## **Bet Constraints**

```proto
// MarketBetConstraints parent group for a market
message MarketBetConstraints {
  string min_amount = 2[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];


  cosmos.base.v1beta1.Coin bet_fee = 3
  [(gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin", (gogoproto.nullable) = false];
}
```

**MarketBetConstraints**: Optional field which can put restrictions for bet acceptance criteria, this optional config help provide more
granular control of different categories of market. We can modify here following:

***MinBetAmount***: Minimum bet amount particular to this market.

***BetFee***: Fixed fee particular to the created market

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
  // bet_constraints holds the constraints of market to accept bets.
  MarketBetConstraints bet_constraints = 9;
  // meta contains human-readable metadata of the market.
  string meta = 10;
  // sr_contribution_for_house is the amount of contribution by sr for the house
  string sr_contribution_for_house = 11 [
    (gogoproto.customname) = "SrContributionForHouse",
    (gogoproto.jsontag) = "sr_contribution_for_house",
    json_name = "sr_contribution_for_house",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // book_uid is the unique identifier corresponding to the book
  string book_uid = 12 [
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

**BetConstraints**: Optional field which can put restrictions for bet acceptance criteria.

**Meta**: Human-Readable data of the market.

**SrContributionForHouse** The contribution of SR for the house corresponding to the market.

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
// Odds is a representation of and market odds items.
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
