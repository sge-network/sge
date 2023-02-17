# **State**


## **Book**

The Book keeps track of the order book for a sport event.

```proto
message Book {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  string sport_event_uid = 1 [(gogoproto.moretags) = "yaml:\"sport_event_uid\""];

  string total_deposit_amount = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"total_deposit_amount\"",
        (gogoproto.nullable) = false];

  string total_bets_amount = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"total_bets_taken\"",
        (gogoproto.nullable) = false];

  string max_loss = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"max_loss\"",
        (gogoproto.nullable) = false];
}
```

## **BookParticipant**

The Book Participant keeps track of the particiapnts of an order book.

```proto
message BookParticipant {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  string sport_event_uid = 1 [(gogoproto.moretags) = "yaml:\"sport_event_uid\""];

  string total_deposit_amount = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"total_deposit_amount\"",
        (gogoproto.nullable) = false];

  string total_bets_amount = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"total_bets_taken\"",
        (gogoproto.nullable) = false];

  string max_loss = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"max_loss\"",
        (gogoproto.nullable) = false];
}
```
