# **State**

## **KVStore**

State in Reward module is defined by its KVStore. This KVStore has one prefix:

1. All campaigns including active and expired items.

The Reward model in the Proto files is as below:

## **Campaign**

1. `creator`: is the creator account(message signer of the CreateCampaign).
2. `uid`: is the universal unique identifier of the campaign.
3. `promoter`: The account address that is responsible for paying the campaign pool balance.
4. `start_ts`: The time that campaign would be started and receive apply reward message.
5. `end_ts`: The time that campaign would be ended and is not able tor receive apply reward message.
6. `reward_category`: Defines the general category of reward that is defined for the campaign.
7. `reward_type`: Defines the type of reward that is defined for the campaign.
8. `reward_amount_type`: Defines the type of reward amount allocation that is defined for the campaign.
9. `reward_amount`: Defines the amount of reward that is defined for the campaign to be granted to main or sub account.
10. `pool`: Information of the current pool balance.
11. `is_active`: Is active/inactive status of the campaign.
12. `claims_per_category`: Maximum number of reward grant transaction per category.
13. `meta`: Contains a string metadata that can be a simple description or a json.

```proto
// Campaign is type for defining the campaign properties.
message Campaign {
  // creator is the address of campaign creator.
  string creator = 1;

  // uid is the unique identifier of a campaign.
  string uid = 2 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // promoter is the address of campaign promoter.
  // Funds for the campaign would be deducted from this account.
  string promoter = 3;

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

  // reward_category is the category of reward.
  RewardCategory reward_category = 6;

  // reward_type is the type of reward.
  RewardType reward_type = 7;

  // amount_type is the type of reward amount.
  RewardAmountType reward_amount_type = 8;

  // reward_amount is the amount defined for a reward.
  RewardAmount reward_amount = 9;

  // pool is the tracker of campaign funds.
  Pool pool = 10 [ (gogoproto.nullable) = false ];

  // is_active is the flag to check if the campaign is active or not.
  bool is_active = 11;

  // claims_per_category is the number of times a user can claim a
  // reward for category of this campaign.
  uint64 claims_per_category = 12;

  // meta is the metadata of the campaign.
  // It is a stringified base64 encoded json.
  string meta = 13;
}

// Pool tracks funds assigned and spent to/for a campaign.
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

// RewardType defines supported types of rewards in reward module.
enum RewardCategory {
  // the invalid or unknown
  REWARD_CATEGORY_UNSPECIFIED = 0;

  // signup reward
  REWARD_CATEGORY_SIGNUP = 1;

  // referral reward
  REWARD_CATEGORY_REFERRAL = 2;

  // affiliate reward
  REWARD_CATEGORY_AFFILIATE = 3;

  // bet refunds
  REWARD_CATEGORY_BET_REFUND = 4;

  // milestone reward
  REWARD_CATEGORY_MILESTONE = 5;

  // bet discounts
  REWARD_CATEGORY_BET_DISCOUNT = 6;

  // other rewards
  REWARD_CATEGORY_OTHER = 7;
}

// RewardType defines supported types of rewards of reward module.
enum RewardType {
  // the invalid or unknown
  REWARD_TYPE_UNSPECIFIED = 0;

  // signup reward
  REWARD_TYPE_SIGNUP = 1;

  // referral signup reward
  REWARD_TYPE_REFERRAL_SIGNUP = 2;

  // affiliate signup reward
  REWARD_TYPE_AFFILIATE_SIGNUP = 3;

  // referral reward
  REWARD_TYPE_REFERRAL = 4;

  // affiliate reward
  REWARD_TYPE_AFFILIATE = 5;

  // bet refunds
  REWARD_TYPE_BET_REFUND = 6;

  // milestone reward
  REWARD_TYPE_MILESTONE = 7;

  // bet discounts
  REWARD_TYPE_BET_DISCOUNT = 8;

  // other rewards
  REWARD_TYPE_OTHER = 9;
}

// RewardType defines supported types of rewards of reward module.
enum RewardAmountType {
  // the invalid or unknown
  REWARD_AMOUNT_TYPE_UNSPECIFIED = 0;

  // Fixed amount
  REWARD_AMOUNT_TYPE_FIXED = 1;

  // Business logic defined amount
  REWARD_AMOUNT_TYPE_LOGIC = 2;

  // Percentage of bet amount
  REWARD_AMOUNT_TYPE_PERCENTAGE = 3;
}
```

## **Reward**

1. `uid`: is the unique universal identifier of the granted reward.
2. `creator`: is the creator account(message signer of the CreateCampaign).
3. `receiver`: is the string address of the main account.
4. `campaign_uid`: is the unique identifier of the associated campaign.
5. `reward_amount`: is the amount to be deducted from main and sub account balances.
6. `source_uid`: is the source of reward grant universal unique identifier.
7. `meta`: is the metadata related to the granted reward, can be a string or json.

```proto
// Reward is the type for transaction made to reward a user
// based on users eligibility.
message Reward {

  // uid is the unique identifier for a reward.
  string uid = 1 [
    (gogoproto.customname) = "UID",
    (gogoproto.jsontag) = "uid",
    json_name = "uid"
  ];

  // creator is the address of the account that invokes the reward transaction.
  string creator = 2;

  // receiver is the address of the account that receives the reward.
  string receiver = 3;

  // campaign_uid is the unique identifier of the campaign.
  string campaign_uid = 4 [
    (gogoproto.customname) = "CampaignUID",
    (gogoproto.jsontag) = "campaign_uid",
    json_name = "campaign_uid"
  ];

  // reward_amount is the amount of the reward.
  RewardAmount reward_amount = 7 [
    (gogoproto.customname) = "RewardAmount",
    (gogoproto.jsontag) = "reward_amount",
    json_name = "reward_amount"
  ];

  // source_uid is the address of the source.
  // It is used to identify the source of the reward.
  // For example, the source uid of a referral signup
  // reward is the address of the referer.
  string source_uid = 8 [
    (gogoproto.customname) = "SourceUID",
    (gogoproto.jsontag) = "source_uid",
    json_name = "source_uid"
  ];

  // meta is the metadata of the campaign.
  // It is a stringified base64 encoded json.
  string meta = 12;
}

// RewardAmount
message RewardAmount {
  // main_account_reward amount transferred to main account address
  string main_account_amount = 1 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"main_account_amount\""
  ];

  // sub_account reward amount transferred to subaccount address
  string subaccount_amount = 2 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"subaccount_amount\""
  ];

  // unlock_period is the period after which the reward is unlocked.
  uint64 unlock_period = 3 [
    (gogoproto.customname) = "UnlockPeriod",
    (gogoproto.jsontag) = "unlock_period",
    json_name = "unlock_period"
  ];
}
```
