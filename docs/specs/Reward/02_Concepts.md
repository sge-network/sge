# **Concepts**

The **Reward Module** is responsible for creating and updating campaigns, as well as managing reward allocation. Users can create campaigns either through the command-line interface or by broadcasting a campaign creation message.

## Campaigns

Campaigns are uniquely identified by their UID. It's possible to have multiple campaigns with the same reward type running concurrently. The application of rewards depends on the specific campaign being defined.

## Reward Types

### SignUp Rewards

1. **SignUp**
   - This reward is granted when a user creates an account by signing up in the system. The reward is stored in the `subaccount` balance and can be used for betting or other house functionalities.

2. **Referral**
   - Users receive this reward when they are referred to the system by another user. Similar to the SignUp reward, the referee's reward is also in the `subaccount` balance and can be utilized for betting or other purposes.

### Referral Rewards

- Referral rewards are given to users who refer new users to the system. The referrer's reward is stored in the `subaccount` balance and can be used for betting or other house functionalities.
