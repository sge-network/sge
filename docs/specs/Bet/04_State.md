# **State**

## **KVStore**

State in bet module is defined by its KVStore. This KVStore has three prefixes:

1. Store for all bets of a certain creator, using this pattern, blockchain is able to return list of all bets, bets of a certain creator and a single bet. The key prefix is created dynamically using this combination: `BetListPrefix`+`{Creator Address}`+`{Secuancial ID}`

2. Store for the map of BetID and BetUID. this helps to get corresponding ID of the bet by issuing the UID. The keys are UIDs of bets and the values are sequencial generated IDs by blockchain.
3. Store for the Bet statistics that contains the count of the total bets used to create next sequencial BetID.

The bet model in the Proto files is as below:

## **Params**

1. `batch_settlement_count`: is the count of bets to be automatically settlement in end-blocker.

```proto
// Params defines the parameters for the module.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // batch_settlement_count is the batch settlement bet counts.
  uint32 batch_settlement_count = 1;
}
```

```proto
// Bet is the main type of bet in the blockchain state.
message Bet {

  // uid is the universal unique identifier assigned to a bet.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // sport_event_uid is the universal unique identifier of
  // the sport-event on which the bet is placed.
  string sport_event_uid = 2 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];

  // odds_uid is the unique universal unique identifier,
  // of the odds on which the bet is placed.
  string odds_uid = 3 [
    (gogoproto.customname) = "OddsUID",
    (gogoproto.jsontag) = "odds_uid",
    json_name = "odds_uid"
  ];

  // odds_type is the type of odds that
  // user choose such as decimal, fractional.
  sgenetwork.sge.bet.OddsType odds_type = 4;

  // odds_value is the odds on which the bet is placed.
  string odds_value = 5;

  // amount is the wager amount.
  string amount = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // bet_fee is the betting fee calculated by the bet amount.
  string bet_fee = 7 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // status is the status of the bet, such as `unspecified` or `settled`.
  Status status = 8;

  // result is the result of the bet, such as `won` or `lost`.
  Result result = 9;

  // verified shows bet is verified or not.
  bool verified = 10;

  // ticket is a signed string containing important info such as `odds_value`.
  string ticket = 11;

  // creator is the bettor address.
  string creator = 12;

  // created_at is the bet placement timestamp.
  int64 created_at = 13;

  // settlement_height is the block height that the bet is settled.
  int64 settlement_height = 14;

  // max_loss_multiplier is the multiplier coefficient of max loss.
  string max_loss_multiplier = 15 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];

  // bet_fulfillment is the fulfillment data.
  repeated BetFulfillment bet_fulfillment = 16;

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

message UID2ID {

  // uid is the unique uuid assigned to bet
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];

  // id is the unique uuid assigned to bet
  uint64 id = 2 [(gogoproto.customname) = "ID" ,(gogoproto.jsontag) = "id", json_name = "id"];

}

message BetStats {

  // count is the total count of bets
  uint64 count = 1;

}

// ActiveBet is the type for an active bet
message ActiveBet {
  // uid is the universal unique identifier for the bet.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  // creator is the bettor address.
  string creator = 2;
}

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

// BetFulfillment is the type for bet fulfillment
message BetFulfillment {
  // participant_address is the bech32-encoded address of the participant
  // fulfilling bet.
  string participant_address = 1
      [ (gogoproto.moretags) = "yaml:\"participant_address\"" ];

  // participation_index is the index in initial participation queue number
  uint64 participation_index = 2
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // bet amount fulfilled by the participation
  string bet_amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"bet_amount\""
  ];

  // payout amount fulfilled by the participation
  string payout_amount = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"payout_amount\""
  ];
}
```
