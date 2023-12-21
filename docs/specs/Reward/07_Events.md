# **Events**

The Bet module emits the following events

## *MsgCreateCampaign*

| **Type**                   | **Attribute Key**         | **Attribute Value**   |
|----------------------------|---------------------------|-----------------------|
| create_campaign            | creator                   | {creator}             |
| create_campaign            | uid                       | {uid}                 |
| message                    | module                    | reward                |
| message                    | action                    | create_campaign       |
| message                    | sender                    | {creator}             |

---

## *MsgUpdateCampaign*

| **Type**                   | **Attribute Key**         | **Attribute Value**   |
|----------------------------|---------------------------|-----------------------|
| update_campaign            | creator                   | {creator}             |
| update_campaign            | uid                       | {uid}                 |
| message                    | module                    | reward                |
| message                    | action                    | update_campaign       |
| message                    | sender                    | {creator}             |

---

## *MsgWithdrawFunds*

| **Type**                   | **Attribute Key**         | **Attribute Value**   |
|----------------------------|---------------------------|-----------------------|
| withdraw_funds             | creator                   | {creator}             |
| withdraw_funds             | uid                       | {uid}                 |
| message                    | module                    | reward                |
| message                    | action                    | withdraw_funds        |
| message                    | sender                    | {creator}             |

---

## *MsgGrantReward*

| **Type**                   | **Attribute Key**         | **Attribute Value**              |
|----------------------------|---------------------------|----------------------------------|
| apply_reward               | creator                   | {creator}                        |
| apply_reward               | campaign_uid              | {campaign_uid}                   |
| apply_reward               | reward_uid                | {reward_uid}                     |
| apply_reward               | promoter                  | {promoter}                       |
| apply_reward               | main_acc_amount           | {main_acc_amount}                |
| apply_reward               | sub_acc_amount            | {sub_acc_amount}                 |
| apply_reward               | sub_acc_unlock_ts         | {sub_acc_unlock_ts}              |
| message                    | module                    | reward                           |
| message                    | action                    | apply_reward                     |
| message                    | sender                    | {creator}                        |

---
