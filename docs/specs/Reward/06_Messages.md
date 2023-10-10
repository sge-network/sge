# **Messages**

In this section, we describe the processing of the Reward messages. the transaction message
handler endpoints is as follows

```proto
service Msg {
  // CreateCampaign is campaign create message endpoint.
  rpc CreateCampaign(MsgCreateCampaign) returns (MsgCreateCampaignResponse);
  // UpdateCampaign is campaign update message endpoint.
  rpc UpdateCampaign(MsgUpdateCampaign) returns (MsgUpdateCampaignResponse);
  // ApplyReward is reward application message endpoint.
  rpc ApplyReward(MsgApplyReward) returns (MsgApplyRewardResponse);
}
```

## **MsgCreateCampaign**

Within this message, the user specifies the campaign information they wish to create.

```proto
// MsgCreateCampaign is campaign create message type.
message MsgCreateCampaign {
  // creator is the address of creator account.
  string creator = 1;
  // uid is the uinque identifier of the campaign.
  string uid = 2;
  // ticket is the payload data.
  string ticket = 3;
}

// MsgCreateCampaignResponse campaign create message response type.
message MsgCreateCampaignResponse {}
```

### **Sample Create Campaign ticket**

```json
{
 "funder_address": "sge1rk85ptmy3gkphlp6wyvuee3lyvz88q6x59jelc",
 "start_ts": 1695970800,
 "end_ts": 1727593200,
 "type": 3,
    "reward_defs": [
        {
            "rec_type": 1,
            "amount": "100",
            "rec_acc_type": 1,
            "unlock_ts": 0
        }
    ],
 "pool_amount": "1000",
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "CreateCampaign"
}
```

## **MsgApplyReward**

Within this message, the user specifies the reward application information they wish to apply.

```proto
// MsgApplyReward is apply reward message type.
message MsgApplyReward {
  // creator is the address of creator account.
  string creator = 1;
  // campaign_uid is the uinque identifier of the campaign.
  string campaign_uid = 2;
  // ticket is the payload data.
  string ticket = 3;
}

// MsgApplyRewardResponse apply reward message response type.
message MsgApplyRewardResponse {}
```

### **Sample Apply Reward ticket**

Note: Signup, Affiliation, and noloss bets rewards needs this format of payload.

```json
{
 "receiver": {
    "rec_type": 1,
    "addr": "sge1w77wnncp6w6llqt0ysgahpxjscg8wspw43jvtd"
 },
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "ApplyReward"
}
```

Note: Referral reward needs this format of payload.

```json
{
 "receivers": [
  {
    "rec_type": 2,
    "addr": "sge1w77wnncp6w6llqt0ysgahpxjscg8wspw43jvtd"
  },
  {
    "rec_type": 3,
    "addr": "sge1afdqdea8r2uh0ujn8l62fw7plvagzqgcmph40n"
  }
 ],
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "ApplyReward"
}
```

Note: All of the definitions of the campaign `reward_defs` should be defined in the `receiver`/`receivers`.
