# **State**

## **KVStore**

State in bet module is defined by its KVStore. This KVStore has three prefixes:

1. Store for all bets of a certain creator, using this pattern, blockchain is able to return list of all bets, bets of a certain creator and a single bet. The key prefix is created dynamically using this combination: `BetListPrefix`+`{Creator Address}`+`{Secuancial ID}`

2. Store for the map of BetID and BetUID. this helps to get corresponding ID of the bet by issuing the UID. The keys are UIDs of bets and the values are sequencial generated IDs by blockchain.
3. Store for the Bet statistics that contains the count of the total bets used to create next sequencial BetID.

The bet model in the Proto files is as below:

```proto
message Bet {

  // uid is the unique uuid assigned to bet
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];

  // sport_event_uid is the unique uuid of the sportevent on which bet is placed
  string sport_event_uid = 2 [(gogoproto.customname) = "SportEventUID" ,(gogoproto.jsontag) = "sport_event_uid", json_name = "sport_event_uid"];

  // odds_uid is the unique uuid of the odds on which bet is placed
  string odds_uid = 3 [(gogoproto.customname) = "OddsUID" ,(gogoproto.jsontag) = "odds_uid", json_name = "odds_uid"];

  // odds_value is the odds on which bet is placed
  string odds_value = 4[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable)   = false];

  // amount is the wager amount
  string amount = 5[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  // betFee is the betting fee
  cosmos.base.v1beta1.Coin bet_fee = 6[
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin",
    (gogoproto.nullable) = false];


  // status is the status of the bet, such as `pending` or `settled`
  Status status = 7;

  // result is the result of bet, such as `won` or `lost`
  Result result = 8;

  // verified shows bet is verified or not
  bool verified = 9;

  // ticket is a signed string containing important info such as `oddsValue`
  string ticket = 10;

  // creator is the bettor address
  string creator = 11;

  // created_at is bet placement timestamp
  int64 created_at = 12;

  // odds_type is the type of odds user chose such as decimal, fractional
  sgenetwork.sge.bet.OddsType odds_type = 13;

  // Status of the Bet.
  enum Status {

    // the unknown status
    STATUS_INVALID = 0;

    // bet is placed
    STATUS_PLACED = 1;

    // bet is canceled by Bettor
    STATUS_CANCELLED = 2;

    // bet is aborted
    STATUS_ABORTED = 3;

    // pending for getting placed
    STATUS_PENDING = 4;

    // bet result is declared
    STATUS_RESULT_DECLARED = 5;

    // the bet is settled
    STATUS_SETTLED = 6;
  }

  // Result of the bet.
  enum Result {

    // the invalid or unknown
    RESULT_INVALID = 0;

    // bet result is pending
    RESULT_PENDING = 1;

    // bet won by the bettor
    RESULT_WON = 2;

    // bet lost by the bettor
    RESULT_LOST = 3;

    // bet is draw
    RESULT_DRAW = 4;

    // bet is aborted
    RESULT_ABORTED = 5;
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

```
