# **Events**

The Sport-Event module emits the following events:

## *MsgEventCreated*

| **Type**                   | **Attribute Key**         | **Attribute Value**   |
|----------------------------|---------------------------|-----------------------|
| create_sport_events        | sport_events_success_uid  | {uid}                 |
| create_sport_events        | sport_events_book_uid     | {book_uid}            |
| message                    | module                    | SportEvent            |
| message                    | action                    | create_sport_events   |
| message                    | sender                    | {creator}             |

---

## *MsgEventResolution*

|   **Type**               |     **Attribute Key**       | **Attribute Value**   |
|:------------------------:|:---------------------------:|:---------------------:|
| update_sport_events      | sport_events_success_uid    | {uid}                 |
| message                  | module                      | SportEvent            |
| message                  | action                      | update_sport_events   |
| message                  | sender                      | {creator}             |

---

## *MsgEventUpdate*

| **Type**                  | **Attribute Key**        | **Attribute Value**   |
|---------------------------|--------------------------|-----------------------|
| resolve_sport_events      | sport_events_success_uid | {uid}                 |
| message                   | module                   | SportEvent            |
| message                   | action                   | resolve_sport_events  |
| message                   | sender                   | {creator}             |
