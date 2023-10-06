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

## *MsgApplyReward*

| **Type**                   | **Attribute Key**         | **Attribute Value**              |
|----------------------------|---------------------------|----------------------------------|
| apply_reward               | creator                   | {creator}                        |
| apply_reward               | campaign_uid              | {campaign_uid}                   |
| distributions              | distributions             | {yaml string of distributions}   |
| message                    | module                    | reward                           |
| message                    | action                    | apply_reward                     |
| message                    | sender                    | {creator}                        |

---
