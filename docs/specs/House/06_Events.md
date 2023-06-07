# **Events**

The House module emits the following events

## *MsgDeposit*

|  Type         |    Attribute Key    |    Attribute Value    |
|:-------------:|:-------------------:|:---------------------:|
| house_deposit | creator             |  {creator}            |
| house_deposit | depositor           |  {depositor}          |
| house_deposit | market_uid          |  {market_uid}         |
| house_deposit | participation_index |  {participation_index}|
| message       | module              |  house                |
| message       | action              |  house_deposit        |
| message       | sender              |  {creator}            |

## *MsgWithdraw*

|  Type          |    Attribute Key    |    Attribute Value    |
|:--------------:|:-------------------:|:---------------------:|
| house_withdraw | creator             |  {creator}            |
| house_withdraw | depositor           |  {depositor}          |
| house_withdraw | market_uid          |  {market_uid}         |
| house_withdraw | participation_index |  {participation_index}|
| message        | module              |  house                |
| message        | action              |  house_withdraw       |
| message        | sender              |  {creator}            |
