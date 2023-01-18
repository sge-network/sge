# **State**

## **Params**

1. `committee_members`: This field will be used to store the list of eligible addresses who can influence the SR via voting. (To be implemented in later phases)

```proto
// Params defines the parameters for the Strategic Reserve module.
message Params{
    option (gogoproto.equal) = true;
    option (gogoproto.goproto_stringer) = false;

    repeated string committee_members = 1
        [(gogoproto.moretags) = "yaml:\"committee_members\""];
}
```

---

## **Reserver**

The reserver keeps track of the current state of the Strategic reserve, that is the amount of tokens available for accepting bets

```proto
// Reserver defines the parameters for the StrategicReserve module.
message Reserver {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  SRPool sr_pool = 1 [(gogoproto.moretags) = "yaml:\"sr_pool\""];
}
```

```proto
// SRPool defines the locked amount and the unlocked amount in the SR Pool Account.
message SRPool {
  option (gogoproto.equal) = true;
  option (gogoproto.goproto_stringer) = false;

  string locked_amount = 1 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"locked_amount\"",
        (gogoproto.nullable) = false];

  string unlocked_amount = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
        (gogoproto.moretags) = "yaml:\"unlocked_amount\"",
        (gogoproto.nullable) = false];
}
```
