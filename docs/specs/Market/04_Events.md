# **Events**

The Market module emits the following events:

## *MsgAddMarket*

| **Type**                   | **Attribute Key**         | **Attribute Value**   |
|----------------------------|---------------------------|-----------------------|
| market_add                 | uid                       | {uid}                 |
| market_add                 | orderbook_uid             | {orderbook_uid}       |
| message                    | module                    | market                |
| message                    | action                    | market_add            |
| message                    | sender                    | {creator}             |

---

## *MsgUpdateMarket*

|   **Type**               |     **Attribute Key**       | **Attribute Value**   |
|:------------------------:|:---------------------------:|:---------------------:|
| market_update            | uid                         | {uid}                 |
| message                  | module                      | market                |
| message                  | action                      | market_update         |
| message                  | sender                      | {creator}             |

---

## *MsgResolveMarket*

| **Type**                  | **Attribute Key**        | **Attribute Value**   |
|---------------------------|--------------------------|-----------------------|
| market_resolve            | uid                      | {uid}                 |
| message                   | module                   | market                |
| message                   | action                   | market_resolve        |
| message                   | sender                   | {creator}             |
