# **State**

## **Params**

1. `min_deposit`: is the minimum allowed amount of deposit.
2. `house_participation_fee`: os the percentage of deposit amount to be paid as deposit fee.

```proto
// Params define the parameters for the house module.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // min_deposit is the minimum amount of acceptable deposit.
  string min_deposit = 1 [
    (gogoproto.moretags) = "yaml:\"min_deposit\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // house_participation_fee is the % of the deposit to be paid for a house
  // participation by the user
  string house_participation_fee = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

---

## **Deposit**

The deposit keeps track of the deposits made byt the users.

```proto
// Deposit represents the deposit against a sport event held by an account.
message Deposit {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // creator is the bech32-encoded address of the depositor.
  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];

  // sport_event_uid is the uid of sport event against which deposit is being
  // made.
  string sport_event_uid = 2 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];

  // participation_index index corresponding to the book participation
  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // amount is the amount being deposited.
  string amount = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"amount\""
  ];

  // fee is deducted from the amount on deposition.
  string fee = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"fee\""
  ];

  // liquidity is the liquidity being provided to the house after fee deduction.
  string liquidity = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"liquidity\""
  ];

  // withdrawal_count is the total count of the withdrawal attempts
  uint64 withdrawal_count = 7 [ (gogoproto.moretags) = "yaml:\"withdrawals\"" ];

  // total_withdrawal_amount is the total amount of withdrawal attempts
  string total_withdrawal_amount = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"total_withdrawal_amount\""
  ];
}
```

---

## **Withdeawal**

The withdrawal keeps track of the withdrawals made byt the users.

```proto
// Withdrawal represents the withdrawal against a deposit.
message Withdrawal {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // withdrawal is the withdrawal attempt id.
  uint64 id = 1 [
    (gogoproto.customname) = "ID",
    (gogoproto.jsontag) = "id",
    json_name = "id",
    (gogoproto.moretags) = "yaml:\"id\""
  ];

  // creator is the bech32-encoded address of the depositor.
  string creator = 2 [ (gogoproto.moretags) = "yaml:\"creator\"" ];

  // sport_event_uid is the uid of sport-event against which the deposit is
  // being made.
  string sport_event_uid = 3 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];

  // participation_index is the id corresponding to the book participation
  uint64 participation_index = 4
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // mode is the withdrawal mode enum value
  WithdrawalMode mode = 5 [ (gogoproto.moretags) = "yaml:\"mode\"" ];

  // amount is the amount being withdrawn.
  string amount = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"amount\""
  ];
}

// WithdrawalMode is the enum type for the withdrawal mode.
enum WithdrawalMode {
  // invalid
  WITHDRAWAL_MODE_UNSPECIFIED = 0;
  // full
  WITHDRAWAL_MODE_FULL = 1;
  // partial
  WITHDRAWAL_MODE_PARTIAL = 2;
}

```
