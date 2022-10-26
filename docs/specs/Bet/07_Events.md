# **Events**

The Strategic Reserve module emits the following events

## *MsgPlaceBet*

|  Type  |  Attribute Key  |  Attribute Value  |
|:-------------:|:---------------:|:---------------:|
| place_bet 	| creator 	|           	|
| place_bet 	| bet_uid  	|             	|
| place_bet 	| sport_event_uid |         |
| place_bet 	| amount  	|             	|
| message      	| module  	|  bet      	|
| message       | action   	| place_bet     |
| message      	| sender    |         	    |

---

## *MsgSettleBet*

|  Type  |  Attribute Key  |  Attribute Value  |
|:-------------:|:---------------:|:---------------:|
| settle_bet 	| bet_creator|           	|
| settle_bet 	| bet_uid  	 |             	|
| settle_bet 	| sport_event_uid |         |
| message      	| module  	 |  bet      	|
| message       | action   	 | settle_bet   |
| message      	| sender     |         	    |

---
