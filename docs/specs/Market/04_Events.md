# **Events**

The Market module emits the following events:

## *MsgMarketCreated*

| **Type**                   | **Attribute Key**         | **Attribute Value**   |
|----------------------------|---------------------------|-----------------------|
| create_markets             | markets_success_uid       | {uid}                 |
| create_markets             | markets_book_uid          | {book_uid}            |
| message                    | module                    | Market                |
| message                    | action                    | create_markets        |
| message                    | sender                    | {creator}             |

---

## *MsgMarketResolution*

|   **Type**               |     **Attribute Key**       | **Attribute Value**   |
|:------------------------:|:---------------------------:|:---------------------:|
| update_markets           | markets_success_uid         | {uid}                 |
| message                  | module                      | Market                |
| message                  | action                      | update_markets        |
| message                  | sender                      | {creator}             |

---

## *MsgMarketUpdate*

| **Type**                  | **Attribute Key**        | **Attribute Value**   |
|---------------------------|--------------------------|-----------------------|
| resolve_markets           | markets_success_uid      | {uid}                 |
| message                   | module                   | Market                |
| message                   | action                   | resolve_markets       |
| message                   | sender                   | {creator}             |
