# **State**

## **Params**

1. `max_book_participations`: is the maximum participations allowed for a book.
2. `batch_settlement_count`: is the count of bets to be automatically settlement in strategic reserve.

```proto
// Params defines the parameters for the strategic reserve module.
message Params {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  // max_order_book_participations is the maximum number of participations per
  // book.
  uint64 max_order_book_participations = 1
      [ (gogoproto.moretags) = "yaml:\"max_order_book_participations\"" ];

  // batch_settlement_count is the batch settlement deposit count.
  uint64 batch_settlement_count = 2
      [ (gogoproto.moretags) = "yaml:\"batch_settlement_count\"" ];

  // requeue_threshold is the threshold at which a participation is requeued in orderbook.
  uint64 requeue_threshold = 3
  [ (gogoproto.moretags) = "yaml:\"requeue_threshold\"" ];
}
```

## **OrderBook**

The OrderBook keeps track of the order book for a market.

```proto
// OrderBook represents the order book maintained against a market.
message OrderBook {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // uid is the universal unique identifier of the order book.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // participation_count is the count of participations in the order book.
  uint64 participation_count = 2
      [ (gogoproto.moretags) = "yaml:\"participation_count\"" ];

  // odds_count is the count of the odds in the order book.
  uint64 odds_count = 3 [ (gogoproto.moretags) = "yaml:\"odds_count\"" ];

  // status represents the status of the order book.
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
// OrderBookParticipation represents the participants of an order book.
message OrderBookParticipation {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // index is the index of the participation in the participation queue.
  uint64 index = 1 [ (gogoproto.moretags) = "yaml:\"index\"" ];

  // order_book_uid is the unique identifier corresponding to the order book.
  string order_book_uid = 2 [
    (gogoproto.customname) = "OrderBookUID",
    (gogoproto.jsontag) = "order_book_uid",
    json_name = "order_book_uid"
  ];

  // participant_address is the bech32-encoded address of the participant.
  string participant_address = 3
      [ (gogoproto.moretags) = "yaml:\"participant_address\"" ];

  // is_module_account represents if the participant is a module account or not.
  bool is_module_account = 4
      [ (gogoproto.moretags) = "yaml:\"is_module_account\"" ];

  // liquidity is the total initial liquidity provided.
  string liquidity = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"liquidity\""
  ];

  // current_round_liquidity is the liquidity provided for the current round.
  string current_round_liquidity = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"current_round_liquidity\""
  ];

  // exposures_not_filled reresents if all of the exposures of the participation
  // are filled or not.
  uint64 exposures_not_filled = 7
      [ (gogoproto.moretags) = "yaml:\"exposures_not_filled\"" ];

  // total_bet_amount is the total bet amount corresponding to all exposures.
  string total_bet_amount = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"total_bet_amount\""
  ];

  // current_round_total_bet_amount is the total bet amount corresponding to all
  // exposures in the current round.
  string current_round_total_bet_amount = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"current_round_total_bet_amount\""
  ];

  // max_loss is the total bet amount corresponding to all exposure.
  string max_loss = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"max_loss\""
  ];

  // current_round_max_loss is the current round max loss.
  string current_round_max_loss = 11 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"current_round_max_loss\""
  ];

  // current_round_max_loss_odds_uid is the total max loss corresponding to
  // all exposures.
  string current_round_max_loss_odds_uid = 12 [
    (gogoproto.customname) = "CurrentRoundMaxLossOddsUID",
    (gogoproto.jsontag) = "current_round_max_loss_odds_uid",
    json_name = "current_round_max_loss_odds_uid",
    (gogoproto.moretags) = "yaml:\"current_round_max_loss_odds_uid\""
  ];

  // actual_profit is the actual profit of the participation fulfillment.
  string actual_profit = 13 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"actual_profit\""
  ];

  // is_settled represents if the participation is settled or not.
  bool is_settled = 14 [ (gogoproto.moretags) = "yaml:\"is_settled\"" ];
}
```

## **ParticipationBetPair**

Keeps track of the bet that is placed and fulfilled.

```proto
// ParticipationBetPair represents the book participation and bet bond.
message ParticipationBetPair {
  // order_book_uid is the unique identifier corresponding to the order book
  string order_book_uid = 1 [
    (gogoproto.customname) = "OrderBookUID",
    (gogoproto.jsontag) = "order_book_uid",
    json_name = "order_book_uid"
  ];

  // participation_index is the index of participation corresponding to the bet
  // fulfillment.
  uint64 participation_index = 2
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // bet_uid is the bet universal unique identifier of the bet that is
  // fulfilled.
  string bet_uid = 3 [
    (gogoproto.customname) = "BetUID",
    (gogoproto.jsontag) = "bet_uid",
    json_name = "bet_uid"
  ];
}
```

## **BetOddsExposure**

Keeps track if Exposures of each odd of order book and market.

```proto
// OrderBookOddsExposure represents the exposures taken on odds.
message OrderBookOddsExposure {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // order_book_uid is the universally unique identifier corresponding to the
  // order book.
  string order_book_uid = 1 [
    (gogoproto.customname) = "OrderBookUID",
    (gogoproto.jsontag) = "order_book_uid",
    json_name = "order_book_uid"
  ];

  // odds_uid is the universally unique identifier of the odds.
  string odds_uid = 2 [
    (gogoproto.customname) = "OddsUID",
    (gogoproto.jsontag) = "odds_uid",
    json_name = "odds_uid"
  ];

  // fulfillment_queue is the slice of indices of participations to be
  // fulfilled.
  repeated uint64 fulfillment_queue = 3
      [ (gogoproto.moretags) = "yaml:\"fulfillment_queue\"" ];
}
```

## **ParticipationExposure**

Keeps track of the exposures of the participation.

```proto
// ParticipationExposure represents the exposures taken on odds by
// participations.
message ParticipationExposure {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // order_book_uid is the universally unique identifier of the order book that
  // the exposure is being set.
  string order_book_uid = 1 [
    (gogoproto.customname) = "OrderBookUID",
    (gogoproto.jsontag) = "order_book_uid",
    json_name = "order_book_uid"
  ];

  // odds_uid is the odds universal unique identifier that the exposure is being
  // set.
  string odds_uid = 2 [
    (gogoproto.customname) = "OddsUID",
    (gogoproto.jsontag) = "odds_uid",
    json_name = "odds_uid"
  ];

  // participation_index is the index of participation in the queue.
  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // exposure is the total exposure taken on given odds.
  string exposure = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"exposure\""
  ];

  // bet_amount is the total bet amount corresponding to the exposure.
  string bet_amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"bet_amount\""
  ];

  // is_fulfilled represents if the participation exposure is fulfilled or not.
  bool is_fulfilled = 6 [ (gogoproto.moretags) = "yaml:\"is_fulfilled\"" ];

  // round is the current round number in the queue.
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
