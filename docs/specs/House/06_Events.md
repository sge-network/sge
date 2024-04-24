# **Events**

The **House module** emits the subsequent events.

## *MsgDeposit*

|  Type         |    Attribute Key    |         Attribute Value           |
|:-------------:|:-------------------:|:---------------------------------:|
| house_deposit | creator             |  {creator}                        |
| house_deposit | depositor           |  {depositor}                      |
| house_deposit | deposit_market_index|  {market_uid#participation_index} |
| message       | module              |  house                            |
| message       | action              |  house_deposit                    |
| message       | sender              |  {creator}                        |

## *MsgWithdraw*

|  Type          |    Attribute Key     |        Attribute Value            |
|:--------------:|:--------------------:|:---------------------------------:|
| house_withdraw | creator              |  {creator}                        |
| house_withdraw | depositor            |  {depositor}                      |
| house_withdraw | withdrawal_id        |  {withdrawal_id}                  |
| house_withdraw | withdraw_market_index|  {market_uid#participation_index} |
| message        | module               |  house                            |
| message        | action               |  house_withdraw                   |
| message        | sender               |  {creator}                        |
