# **State**

## **Params**

1. `max_book_participations`: is the maximum participations allowed for a book.
2. `batch_settlement_count`: is the count of bets to be automatically settlement in order book.

```proto
// Params defines the parameters for the orderbook module.
message Params {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  // max_book_participations is the maximum number of participations per book.
  uint64 max_book_participations = 1
      [ (gogoproto.moretags) = "yaml:\"max_book_participations\"" ];

  // batch_settlement_count is the batch settlement deposit counts.
  uint64 batch_settlement_count = 2
      [ (gogoproto.moretags) = "yaml:\"batch_settlement_count\"" ];
}
```

## **OrderBook**

The OrderBook keeps track of the order book for a sport-event.

```proto
message OrderBook {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // id corresponding to the book
  string id = 1 [
    (gogoproto.customname) = "ID",
    (gogoproto.jsontag) = "id",
    json_name = "id"
  ];

  // participation_count is the count of participations in the order book
  uint64 participation_count = 2
      [ (gogoproto.moretags) = "yaml:\"participation_count\"" ];

  // odds_count is the count of the odds in the order book
  uint64 odds_count = 3 [ (gogoproto.moretags) = "yaml:\"odds_count\"" ];

  // order book status
  OrderBookStatus status = 4;
}

// OrderBookStatus is the enum type for the status of the orderbook.
enum OrderBookStatus {
  // invalid
  ORDER_BOOK_STATUS_UNSPECIFIED = 0;
  // active
  ORDER_BOOK_STATUS_STATUS_ACTIVE = 1;
  // resolved not settled
  ORDER_BOOK_STATUS_STATUS_RESOLVED = 2;
  // resolved and settled
  ORDER_BOOK_STATUS_STATUS_SETTLED = 3;
}
```

## **BookParticipation**

The Book BookParticipation keeps track of the particiapation of an order book.

```proto
// BookParticipation represents the participants of an order book.
message BookParticipation {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // index is the id of initial participation queue
  uint64 index = 1 [ (gogoproto.moretags) = "yaml:\"index\"" ];

  // book id is id corresponding to the book
  string book_id = 2 [
    (gogoproto.customname) = "BookID",
    (gogoproto.jsontag) = "book_id",
    json_name = "book_id"
  ];

  // participant_address is the bech32-encoded address of the participant.
  string participant_address = 3
      [ (gogoproto.moretags) = "yaml:\"participant_address\"" ];

  // if participation is a module account
  bool is_module_account = 4
      [ (gogoproto.moretags) = "yaml:\"is_module_account\"" ];
  ;

  // liquidity is the total initial liquidity provided
  string liquidity = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"liquidity\""
  ];

  // current round liquidity is the liquidity provided for current round
  string current_round_liquidity = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"current_round_liquidity\""
  ];

  uint64 exposures_not_filled = 7
      [ (gogoproto.moretags) = "yaml:\"exposures_not_filled\"" ];
  ;

  // total_bet_amount is the total bet amount corresponding to all exposure
  string total_bet_amount = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"total_bet_amount\""
  ];

  // current_round_total_bet_amount is the total bet amount corresponding to all
  // exposure
  string current_round_total_bet_amount = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"current_round_total_bet_amount\""
  ];

  // max_loss is the total bet amount corresponding to all exposure
  string max_loss = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"max_loss\""
  ];

  // current_round_max_loss is the total bet amount corresponding to all
  // exposure
  string current_round_max_loss = 11 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"current_round_max_loss\""
  ];

  // current_round_max_loss_odds_uid is the total bet amount corresponding to all
  // exposure
  string current_round_max_loss_odds_uid = 12 [
    (gogoproto.customname) = "CurrentRoundMaxLossOddsUID",
    (gogoproto.jsontag) = "current_round_max_loss_odds_uid",
    json_name = "current_round_max_loss_odds_uid",
    (gogoproto.moretags) = "yaml:\"current_round_max_loss_odds_uid\""
  ];

  // actual_profit is the actual profit
  string actual_profit = 13 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"actual_profit\""
  ];

  // if participation is settled
  bool is_settled = 14 [ (gogoproto.moretags) = "yaml:\"is_settled\"" ];
  ;
}
```

## **ParticipationBetPair**

Keeps track of the bet that is placed and fulfilled.

```proto
// ParticipationBetPair represents the book participation and bet bond
message ParticipationBetPair {
  // book id is id corresponding to the book
  string book_id = 1 [
    (gogoproto.customname) = "BookID",
    (gogoproto.jsontag) = "book_id",
    json_name = "book_id"
  ];

  // participation_index is the count of initial participation queue
  uint64 participation_index = 2
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // bet_uid is bet's uuid
  string bet_uid = 3 [
    (gogoproto.customname) = "BetUID",
    (gogoproto.jsontag) = "bet_uid",
    json_name = "bet_uid"
  ];
}
```

## **BetOddsExposure**

Keeps track if Exposures of each odds of order book and sport-event.

```proto
// BookOddsExposure represents the exposures taken on odds
message BookOddsExposure {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // book_id is id corresponding to the book
  string book_id = 1 [
    (gogoproto.customname) = "BookID",
    (gogoproto.jsontag) = "book_id",
    json_name = "book_id"
  ];

  // odds_uid is odd'd uid
  string odds_uid = 2 [
    (gogoproto.customname) = "OddsUID",
    (gogoproto.jsontag) = "odds_uid",
    json_name = "odds_uid"
  ];

  repeated uint64 fulfillment_queue = 3
      [ (gogoproto.moretags) = "yaml:\"fulfillment_queue\"" ];
}
```

## **ParticipationExposure**

Keeps track of the exposures of the a participation.

```proto

// ParticipationExposure represents the exposures taken on odds by
// participations
message ParticipationExposure {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // book_id is id corresponding to the book
  string book_id = 1 [
    (gogoproto.customname) = "BookID",
    (gogoproto.jsontag) = "book_id",
    json_name = "book_id"
  ];

  // odds_uid is odd's uid
  string odds_uid = 2 [
    (gogoproto.customname) = "OddsUID",
    (gogoproto.jsontag) = "odds_uid",
    json_name = "odds_uid"
  ];

  // participation_index is the id of initial participation queue
  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // exposure is the total exposure taken on given odd
  string exposure = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"exposure\""
  ];

  // bet_amount is the total bet amount corresponding to the exposure
  string bet_amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"bet_amount\""
  ];

  bool is_fulfilled = 6 [ (gogoproto.moretags) = "yaml:\"is_fulfilled\"" ];

  // number of current round in queue
  uint64 round = 7 [ (gogoproto.moretags) = "yaml:\"round\"" ];
}
```

## **OrderBookStats**

Keeps track of statistics of the order book.

```proto
// OrderBookStats holds statistics of the order-book
message OrderBookStats {
  // resolved_unsettled is the list of book ids that needs to be settled.
  repeated string resolved_unsettled = 1;
}
```
