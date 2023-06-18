# **State**

## **Params**

1. `min_deposit`: is the minimum allowed amount of deposit.
2. `house_participation_fee`: os the percentage of deposit amount to be paid as deposit fee.

```proto
// Params define the parameters for the house module.
message Params {
  option (gogoproto.goproto_stringer) = false;

  // min_deposit is the minimum acceptable deposit amount.
  string min_deposit = 1 [
    (gogoproto.moretags) = "yaml:\"min_deposit\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];

  // house_participation_fee is the % of the deposit to be paid for a house
  // participation by the depositor.
  string house_participation_fee = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

---

## **Deposit**

The deposit keeps track of the deposits made by the users who want to participate as house.

```proto
// Deposit represents the deposit against a market held by an account.
message Deposit {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // creator is the bech32-encoded address of the depositor.
  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];

  // creator is the bech32-encoded address of the depositor.
  string depositor_address = 2 [ (gogoproto.moretags) = "yaml:\"depositor_address\"" ];

  // market_uid is the uid of market/order book against which deposit is being
  // made.
  string market_uid = 3 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];

  // participation_index is the index corresponding to the order book
  // participation
  uint64 participation_index = 4
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // amount is the amount being deposited on an order book to be a house
  string amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"amount\""
  ];

  // fee is deducted from the deposited amount for participation in the order
  // book.
  string fee = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"fee\""
  ];

  // liquidity is the liquidity being provided to the order book after fee
  // deduction.
  string liquidity = 7 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"liquidity\""
  ];

  // withdrawal_count is the total count of the withdrawals from an order book
  uint64 withdrawal_count = 8 [ (gogoproto.moretags) = "yaml:\"withdrawals\"" ];

  // total_withdrawal_amount is the total amount withdrawn from the liquidity
  // provided
  string total_withdrawal_amount = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"total_withdrawal_amount\""
  ];
}
```

---

## **Withdrawal**

The withdrawal keeps track of the withdrawals made by the depositor accounts.

```proto
// Withdrawal represents the withdrawal against a deposit.
message Withdrawal {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (gogoproto.goproto_stringer) = false;

  // creator is the bech32-encoded address of the depositor.
  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];

  // withdrawal is the withdrawal attempt id.
  uint64 id = 2 [
    (gogoproto.customname) = "ID",
    (gogoproto.jsontag) = "id",
    json_name = "id",
    (gogoproto.moretags) = "yaml:\"id\""
  ];

  // address is the bech32-encoded address of the depositor.
  string address = 3 [ (gogoproto.moretags) = "yaml:\"address\"" ];

  // market_uid is the uid of market against which the deposit is
  // being made.
  string market_uid = 4 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];

  // participation_index is the id corresponding to the book participation
  uint64 participation_index = 5
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];

  // mode is the withdrawal mode enum value
  WithdrawalMode mode = 6 [ (gogoproto.moretags) = "yaml:\"mode\"" ];

  // amount is the amount being withdrawn.
  string amount = 7 [
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
