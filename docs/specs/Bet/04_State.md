# **State**

## **KVStore**

State in bet module is defined by its KVStore. This KVStore has five prefixes:

1. All bets of a certain creator, using this pattern, blockchain is able to return list of all bets, bets of a certain creator and a single bet. The key prefix is created dynamically using this combination: `BetListPrefix`+`{Creator Address}`+`{Secuential Bet ID}`

2. Map of BetID and BetUID. this helps to get corresponding ID of the bet by issuing the UID. The keys are UIDs of bets and the values are sequencial generated IDs by blockchain.
3. Pending bets of a certain Market to help batch settlement.
4. Settled bets of a block height to keep track of the settled bets for the oracle services.
5. Bet statistics that contains the count of the total bets used to create next sequencial BetID.

The bet model in the Proto files is as below:

## **Params**

1. `batch_settlement_count`: is the count of bets to be automatically settlement in end-blocker.
2. `max_bet_by_uid_query_count`: is the max count of bets to be returned in the bets by uids query.
3. `constraints` contains criteria of the bet placement.
    - `min_amount` minimum bet amount while placement.
    - `fee` bet fee amount payable by bettor.

```proto
// Params defines the parameters for the module.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // batch_settlement_count is the batch settlement bet count.
  uint32 batch_settlement_count = 1;
  // max_bet_by_uid_query_count is the maximum bet by uid query items count.
  uint32 max_bet_by_uid_query_count = 2;
  // constraints is the bet constraints.
  Constraints constraints = 3 [
    (gogoproto.moretags) = "yaml:\"constraints\"",
    (gogoproto.nullable) = false
  ];
}
```

## **Constraints**

Holds bet placement constraints of the bet module.

```proto
// Constraints is the bet constrains type for the bets
message Constraints {
  // min_amount is the minimum allowed bet amount.
  string min_amount = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // fee is the fee that the bettor needs to pay to bet.
  string fee = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

## **Bet**

Bet contains the properties of bet object type.

```proto
// Bet is the transaction order placed by a bettor on a specific event and odd
message Bet {
  // uid is the universal unique identifier assigned to a bet.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // market_uid is the universal unique identifier of
  // the market on which the bet is placed.
  string market_uid = 2 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];

  // odds_uid is the universal unique identifier,
  // of the odds on which the bet is placed.
  string odds_uid = 3 [
    (gogoproto.customname) = "OddsUID",
    (gogoproto.jsontag) = "odds_uid",
    json_name = "odds_uid"
  ];

  // odds_value is the odds on which the bet is placed.
  string odds_value = 4;

  // amount is the wager amount.
  string amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // fee is the betting fee user needs to pay for placing a bet
  string fee = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // status is the status of the bet, such as `unspecified` or `settled`.
  Status status = 7;

  // result is the result of the bet, such as `won` or `lost`.
  Result result = 8;

  // creator is the bettor address.
  string creator = 9;

  // created_at is the bet placement timestamp.
  int64 created_at = 10;

  // settlement_height is the block height at which the bet is settled.
  int64 settlement_height = 11;

  // max_loss_multiplier is the multiplier coefficient of max loss.
  string max_loss_multiplier = 12 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];

  // bet_fulfillment is the fulfillment data.
  repeated BetFulfillment bet_fulfillment = 13;

  // meta is metadata for bet
  MetaData meta = 14;

  // Status of the Bet.
  enum Status {
    // the invalid or unknown
    STATUS_UNSPECIFIED = 0;
    // bet is placed
    STATUS_PLACED = 1;
    // bet is canceled by Bettor
    STATUS_CANCELED = 2;
    // bet is aborted
    STATUS_ABORTED = 3;
    // bet is pending for getting placed
    STATUS_PENDING = 4;
    // bet result is declared
    STATUS_RESULT_DECLARED = 5;
    // the bet is settled
    STATUS_SETTLED = 6;
  }

  // Result of the bet.
  enum Result {
    // the invalid or unknown
    RESULT_UNSPECIFIED = 0;
    // bet result is pending
    RESULT_PENDING = 1;
    // bet won by the bettor
    RESULT_WON = 2;
    // bet lost by the bettor
    RESULT_LOST = 3;
    // bet is refunded
    RESULT_REFUNDED = 4;
  }
}
```

## **UID2ID**

Used for storing the map of bet UUID and sequencial generated numeric ID.

```proto
// UID2ID is the type for mapping UIDs and Sequential IDs of bets.
message UID2ID {
  // uid is the universal unique identifier assigned to the bet.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // id is an autogenerated sequential id for a bet.
  uint64 id = 2 [
    (gogoproto.customname) = "ID",
    (gogoproto.jsontag) = "id",
    json_name = "id"
  ];
}
```

## **PendingBet**

A bet that is not settled yet is pending, so blockchain keeps track of it
the UUID list of these kind of bets.

```proto
// PendingBet is the type for an unsettled bet
message PendingBet {
  // uid is the universal unique identifier for the bet.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // creator is the bettor address.
  string creator = 2;
}
```

## **SettledBet**

A bet that is canceled, aborted or its result is declared is stored in this type.

```proto
// SettledBet is the type for a settled bet.
message SettledBet {
  // uid is the universal unique identifier for the bet.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // bettor_address is the bech32 address of the bettor account.
  string bettor_address = 2;
}
```

## **BetFulfillment**

The `orderbook` module's end blocker, process the settled markets and corresponsing
bets, the payout fulfillment is stored in the bet module state for each bet with this type.

```proto
// BetFulfillment: A bet can be fulfilled by multiple users participating as a
// house Every participant is exposed to a share of risk or payout associated
// with the bet For the risk exposure on a bet, an estimated bet amount is also
// allocated to the participant This bet amount is the amount participant
// receive if the bettor loose the bet
message BetFulfillment {
  // participant_address is the bech32-encoded address of the participant
  // fulfilling bet.
  string participant_address = 1
      [ (gogoproto.moretags) = "yaml:\"participant_address\"" ];
  // participation_index is the index in initial participation queue index
  uint64 participation_index = 2
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
  // bet amount fulfilled by the participation
  string bet_amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"bet_amount\""
  ];
  // payout_profit is the fulfilled profit by the participation.
  string payout_profit = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"payout_profit\""
  ];
}
```
