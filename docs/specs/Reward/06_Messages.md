# **Messages**

In this section, we describe the processing of the Reward messages. the transaction message
handler endpoints is as follows

```proto
service Msg {
  // CreateCampaign is a method to create a campaign
  rpc CreateCampaign(MsgCreateCampaign) returns (MsgCreateCampaignResponse);
  // UpdateCampaign is a method to update campaign
  rpc UpdateCampaign(MsgUpdateCampaign) returns (MsgUpdateCampaignResponse);
  // WithdrawCampaignFunds is method to withdraw funds from the campaign
  rpc WithdrawFunds(MsgWithdrawFunds) returns (MsgWithdrawFundsResponse);
  // GrantReward is method to allocate rewards
  rpc GrantReward(MsgGrantReward) returns (MsgGrantRewardResponse);
}
```

## **MsgCreateCampaign**

Within this message, the user specifies the campaign information they wish to create.

```proto
// MsgCreateCampaign is msg to create a reward campaign
message MsgCreateCampaign {
  // creator is the address of campaign creator account.
  string creator = 1;
  // uid is the unique identifier of the campaign.
  string uid = 2;
  // total_funds is the total funds allocated to the campaign.
  string total_funds = 3 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"total_funds\""
  ];
  // ticket is the payload data.
  string ticket = 4;
}

// MsgCreateCampaignResponse campaign create message response type.
message MsgCreateCampaignResponse {}
```

### **Sample Create Campaign ticket**

```json
{
 "promoter": "sge1rk85ptmy3gkphlp6wyvuee3lyvz88q6x59jelc",
 "start_ts": 1695970800,
 "end_ts": 1727593200,
 "category": 1,
 "reward_type": 1,
 "reward_amount_type": 1,
 "reward_amount": {
  "main_account_amount": "100",
  "sub_account_amount": "100",
  "unlock_period": 136318754373
 },
 "is_active": true,
 "claims_per_category": 1,
 "meta": "custom metadata",
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "CreateCampaign"
}
```

## **MsgUpdateCampaign**

Within this message, the user specifies the campaign information they wish to update.

```proto
// MsgUpdateCampaign is campaign update message type.
message MsgUpdateCampaign {
  // creator is the address of creator account.
  string creator = 1;
  // uid is the unique identifier of the campaign.
  string uid = 2;
  // topup_funds is the topup funds to increase the pool balance of the
  // campaign.
  string topup_funds = 3 [
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"topup_funds\""
  ];
  // ticket is the payload data.
  string ticket = 4;
}

// MsgUpdateCampaignResponse campaign update message response type.
message MsgUpdateCampaignResponse {}
```

### **Sample Update Campaign ticket**

```json
{
 "end_ts": 1727593200,
 "is_active": true,
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "UpdateCampaign"
}
```

## **MsgGrantReward**

Within this message, the user specifies the reward application information they wish to apply.

```proto
// MsgGrantReward is execute reward message type.
message MsgGrantReward {
  // creator is the address of creator account.
  string creator = 1;
  // uid is the unique identifier of the reward.
  string uid = 2;
  // campaign_uid is the unique identifier of the reward campaign.
  string campaign_uid = 3;
  // ticket is the payload data.
  string ticket = 4;
}

// MsgApplyRewardResponse apply reward message response type.
message MsgApplyRewardResponse {}
```

### MsgGrantRewardTicket Payloads

#### Common ticket payload type

This is a common type that should be present in all of the ticket payloads

```proto
// RewardPayloadCommon
message RewardPayloadCommon {
  // receiver is the address of the account that receives the reward.
  string receiver = 1;

  // source_uid is the address of the source.
  // It is used to identify the source of the reward.
  // For example, the source uid of a referral signup reward is the address of
  // the referer.
  string source_uid = 2 [
    (gogoproto.customname) = "SourceUID",
    (gogoproto.jsontag) = "source_uid",
    json_name = "source_uid"
  ];

  // meta is the metadata of the campaign.
  // It is a stringified base64 encoded json.
  string meta = 3;
}
```

#### **1. Individual Signup**

```proto
// GrantSignupRewardPayload is the type for signup reward grant payload.
message GrantSignupRewardPayload {
  // common is the common properties of a reward
  RewardPayloadCommon common = 2 [ (gogoproto.nullable) = false ];
}
```

##### **Sample Signup Grant Reward ticket**

The `source_uid` should be empty.

```json
{
 "common": {
  "receiver": "sge1rk85ptmy3gkphlp6wyvuee3lyvz88q6x59jelc",
  "source_uid": "",
  "meta": "custom grant metadata"
 },
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "GrantReward"
}
```

#### **2. Referral Sign Up**

The `source_uid` is the address of the referrer.

```proto
// GrantSignupRewardPayload is the type for signup reward grant payload.
message GrantSignupRewardPayload {
  // common is the common properties of a reward
  RewardPayloadCommon common = 2 [ (gogoproto.nullable) = false ];
}
```

### **Sample Referral Signup Grant Reward ticket**

```json
{
 "common": {
  "receiver": "sge1rk85ptmy3gkphlp6wyvuee3lyvz88q6x59jelc",
  "source_uid": "{account address of the referrer}",
  "meta": "custom grant metadata"
 },
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "GrantReward"
}
```

#### **2. Referral**

- The `source_uid` should be empty.
- The `referee` should contain the address of referee who already claimed their signup referral reward.

```proto
// GrantSignupReferrerRewardPayload is the type for signup referrer reward grant
// payload.
message GrantSignupReferrerRewardPayload {
  // common is the common properties of a reward
  RewardPayloadCommon common = 1 [ (gogoproto.nullable) = false ];

  // referee is the address of the account that used this referrer address as
  // source_uid
  string referee = 2;
}
```

#### **Sample Referral Grant Reward ticket**

```json
{
 "common": {
  "receiver": "sge1rk85ptmy3gkphlp6wyvuee3lyvz88q6x59jelc",
  "source_uid": "",
  "meta": "custom grant metadata"
 },
 "referee": "{account address of the referee who has already claimed their reward}",
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "GrantReward"
}
```

## **MsgWithdrawFunds**

Within this message, the user specifies the reward application information they wish to apply.

```proto
// MsgWithdrawFunds is withdraw funds message type.
message MsgWithdrawFunds {
  // creator is the address of creator account.
  string creator = 1;
  // uid is the unique identifier of the reward campaign.
  string uid = 2;
  // ticket is the payload data.
  string ticket = 3;
}

// MsgWithdrawFundsResponse withdraw funds message response type.
message MsgWithdrawFundsResponse {}
```

### **Sample Withdraw funds ticket**

```json
{
 "promoter": "sge1rk85ptmy3gkphlp6wyvuee3lyvz88q6x59jelc",
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "WithdrawFunds"
}
```
