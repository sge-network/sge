# **Events**

The Sport-Event module emits the following events:

## *MsgEventCreated*

| **Type**                   | **Attribute Key**     | **Attribute Value**   |
|----------------------------|-----------------------|-----------------------|
| complete_event_creation    | event_id              |                       |
| complete_event_creation    | start_ts              |                       |
| complete_event_creation    | end_ts                |                       |
| complete_event_creation    | odds_uids             | [ ] string            |
| complete_event_creation    | status                |                       |
| complete_event_creation    | max_bet_cap           |                       |
| complete_event_creation    | min_amount            |                       |
| complete_event_creation    | bet_fee               |                       |
| complete_event_creation    | current_total_amount  |                       |
| complete_event_creation    | max_loss              |                       |
| complete_event_creation    | max_vig               |                       |
| complete_event_creation    | min_vig               |                       |
| message                    | module                | SportEvent            |
| message                    | action                | event_creation        |
| message                    | sender                |                       |

---

## *MsgEventResolution*

|   **Type**         |     **Attribute Key**       | **Attribute Value**   |
|:------------------:|:---------------------------:|:---------------------:|
| event_resolved     |     event_id                |                       |
| event_resolved     | resolution_ts               |                       |
| event_resolved     |    status                   |                       |

---

## *MsgEventUpdate*

| **Type**                  | **Attribute Key**     | **Attribute Value**   |
|---------------------------|-----------------------|-----------------------|
| event_updated             | event_id              |                       |
| event_updated             | start_ts              |                       |
| event_updated             | end_ts                |                       |
| event_updated             | odds_uids             | [ ] string            |
| event_updated             | status                |                       |
| event_updated             | max_bet_cap           |                       |
| event_updated             | min_amount            |                       |
| event_updated             | bet_fee               |                       |
| event_updated             | current_total_amount  |                       |
| event_updated             | max_loss              |                       |
| event_updated             | max_vig               |                       |
| event_updated             | min_vig               |                       |
| message                   | module                | SportEvent            |
| message                   | action                | event_updated         |
| message                   | sender                |                       |
