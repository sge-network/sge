# **State**

## **KVStore**

The state within the **Subaccount Module** is defined by its KVStore, which consists of five prefixes:

1. **Sub Account Sequential ID Store**: This prefix keeps track of the last generated ID for sub accounts.
2. **Sub Account Owner**: Establishes a one-to-one relationship between the main account and its associated `subaccount`.
3. **Sub Account Address**: Enables the blockchain to retrieve `subaccount` information directly using the `subaccount` address.
4. **Locked Balance of Each Subaccount Address**: Records the locked balance of each `subaccount` at a specific point in time.
5. **Locked Balance of Each Subaccount**: Stores the overall locked balance for each `subaccount`.

## Parameters

1. **Wager Enabled**: Determines whether wagering via `subaccount` is enabled.
2. **Deposit Enabled**: Indicates whether depositing to the house via `subaccount` is allowed.

```proto
// Params defines the parameters for the module.
message Params { 
  option (gogoproto.goproto_stringer) = false;

  // wager_enabled is enable/disable status of wager feature.
  bool wager_enabled = 1;
  // deposit_enabled is enable/disable status of deposit feature.
  bool deposit_enabled = 2;
}
```

### **AccountSummary**

This type contains the current state of the sub account balance.

```proto
// AccountSummary defines the balance of a subaccount.
message AccountSummary {
  // deposited_amount keeps track of how much was deposited so far in the
  // subaccount.
  string deposited_amount = 1 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];

  // spent_amount keeps track of how much was spent in the account in betting,
  // house, staking, etc.
  string spent_amount = 2 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];

  // withdrawn_amount keeps track of how much was withdrawn in the account after
  // locked coins become unlocked.
  string withdrawn_amount = 3 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];

  // lost_amount keeps track of the amounts that were lost due to betting
  // losses, slashing etc.
  string lost_amount = 4 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];
}
```

### **LockedBalance**

This type contains the locked balance state of a sub account.

```proto
// LockedBalance defines a balance which is locked.
message LockedBalance {
  uint64 unlock_ts = 1 [
    (gogoproto.customname) = "UnlockTS",
    (gogoproto.jsontag) = "unlock_ts",
    json_name = "unlock_ts"
  ];
  string amount = 2 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];
}
```
