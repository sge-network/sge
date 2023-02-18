# **State**

## **Params**

1. `max_depositors`: This field will be used to store the maximum number of despositors allowed per sport event.

```proto
// Params defines the parameters for the house module.
message Params {
  option (gogoproto.equal)            = true;
  option (gogoproto.goproto_stringer) = false;

  // max_depositors is the maximum number of depositors per sport event.
  uint32 max_depositors = 1 [(gogoproto.moretags) = "yaml:\"max_depositors\""];
}
```

---

## **deposits**

The deposits keeps track of the deposit made.

```proto
/ Deposit represents the deposit against a sport event held by an account.
message Deposit {
    option (gogoproto.equal)            = false;
    option (gogoproto.goproto_getters)  = false;
    option (gogoproto.goproto_stringer) = false;

  // depositor_address is the bech32-encoded address of the depositor.
  string depositor_address = 1 [(gogoproto.moretags) = "yaml:\"depositor_address\""];

  string sport_event_uid = 2 [(gogoproto.moretags) = "yaml:\"sport_event_uid\""];

  cosmos.base.v1beta1.Coin deposit_amount = 3 [
    (gogoproto.moretags) = "yaml:\"deposit_amount\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Coin",
    (gogoproto.nullable)   = false];

  cosmos.base.v1beta1.Coin fee = 4 [
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin",
    (gogoproto.nullable) = false];

  cosmos.base.v1beta1.Coin withdrawn_amount = 5 [
    (gogoproto.moretags) = "yaml:\"withdrawn_amount\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Coin",
    (gogoproto.nullable)   = false];
}
```
