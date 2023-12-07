# **State**

## **KVStore**

State in sub account module is defined by its KVStore. This KVStore has five prefixes:

1. Sub account sequential ID store to track the last generated ID of the sub accounts.
2. Sub account owner, that makes a 1-1 relation between main account and the subaccount.
3. Sub account address, enable the blockchain to get the subaccount info by subaccount address itself.
4. Locked balance of each subaccount address at a certain time.
5. Locked balance of each subaccount.

### **Balance**

This type contains the current state of the sub account balance.

```proto
// Balance defines the balance of a subaccount.
message Balance {
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
  // locked coins become free.
  string withdrawm_amount = 3 [
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

### **LockerBalance**

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
