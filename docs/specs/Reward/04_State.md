# **State**

## **KVStore**

State in Reward module is defined by its KVStore. This KVStore has one prefix:

1. All campaigns including active and expired items.

The Reward model in the Proto files is as below:

## **Campaign**

1. `creator`: is the creator account(message signer of the CreateCampaign).
2. `uid`: is the universal unique identifier of the campaign.
3. `funder_address`: The account address that is responsible for paying the campaing pool balance.
4. `start_ts`: The time that campaign would be started and receive apply reward message.
5. `end_ts`: The time that campaign would be ended and is not able tor receive apply reward message.
6. `reward_type`: Defines the type of reward that is defined for the campaign.
7. `reward_defs`: Definitions of the payable rewards receivers and amounts.
8. `pool`: Information of the current pool balance.

```proto
// Campaign is type for defining the campaign properties.
message Campaign {
  // creator is the address of campaign creator.
  string creator = 1;
  string uid = 2 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];
  string funder_address = 3;
  // start_ts is the start timestamp of a campaign.
  uint64 start_ts = 4 [
    (gogoproto.customname) = "StartTS",
    (gogoproto.jsontag) = "start_ts",
    json_name = "start_ts"
  ];
  // end_ts is the end timestamp of a campaign.
  uint64 end_ts = 5 [
    (gogoproto.customname) = "EndTS",
    (gogoproto.jsontag) = "end_ts",
    json_name = "end_ts"
  ];
  // reward_type is the type of defined reward.
  RewardType reward_type = 6;
  // reward_defs is the list of definitions of the campaign rewards.
  repeated Definition reward_defs = 7 [ (gogoproto.nullable) = false ];
  // pool is the tracker of pool funds of the campaign.
  Pool pool = 8 [ (gogoproto.nullable) = false ];
}

// Pool is the type for the campaign funding pool.
message Pool {
  string total = 1 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"total\""
  ];
  string spent = 2 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"spent\""
  ];
}

// Definition is the type for reward declaration for a campaign.
message Definition {
  ReceiverType rec_type = 1;
  ReceiverAccType rec_acc_type = 2;
  string amount = 3 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"amount\""
  ];
  uint64 unlock_ts = 4 [
    (gogoproto.customname) = "UnlockTS",
    (gogoproto.jsontag) = "unlock_ts",
    json_name = "unlock_ts"
  ];
}

// Receiver is the type for reward receiver properties.
message Receiver {
  ReceiverType rec_type = 1;
  string addr = 2;
}

// RewardType defines supported types of rewards of reward module.
enum RewardType {
  // the invalid or unknown
  REWARD_TYPE_UNSPECIFIED = 0;
  // signup reward
  REWARD_TYPE_SIGNUP = 1;
  // referral reward
  REWARD_TYPE_REFERRAL = 2;
  // affiliation reward
  REWARD_TYPE_AFFILIATION = 3;
  // noloss bets reward
  REWARD_TYPE_NOLOSS_BETS = 4;
}

// ReceiverAccType defines supported types account types for reward
// distribution.
enum ReceiverAccType {
  // the invalid or unknown
  RECEIVER_ACC_TYPE_UNSPECIFIED = 0;
  // main account
  RECEIVER_ACC_TYPE_MAIN = 1;
  // sub account
  RECEIVER_ACC_TYPE_SUB = 2;
}

// ReceiverType defines all of reward receiver types in the system.
enum ReceiverType {
  // the invalid or unknown
  RECEIVER_TYPE_UNSPECIFIED = 0;
  // single receiver account
  RECEIVER_TYPE_SINGLE = 1;
  // referrer
  RECEIVER_TYPE_REFERRER = 2;
  // referee
  RECEIVER_TYPE_REFEREE = 3;
}
```
