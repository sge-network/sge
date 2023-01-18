# **State Transitions**

This section defines the state transitions of the sport module's KVStore in all scenarios:

## **Add SportEvent**

When this is processed:

- Validate the creator address and validate the ticket format
- Call the DVM module to validate the ticket internals and to retrieve the
  contents of the ticket
- If the ticket is valid, check if the event already exists
- If the sport-event does not already exist, a new sport event should be
  created with the given data and will be added to the `sportevent module's KVStore`. Like this:

```go
newEvent := &type.SportEvent{
 Uid            : <string>
 StartTS        : <uint64>
 EndTS          : <uint64>
 Odds           : <[]*Odds>
 WinnerOddsUIDs : <map[string][]byte>
 Status         : <SportEventStatus>
 ResolutionTS   : <uint64>
 Creator        : <string>
 BetConstraints : <*EventBetConstraints>
 Active         : <bool>
 Meta           : <string>
}
```

---

## **Update SportEvent**

When this is processed:

For each event:

- Validate the creator address and validate the ticket format
- Call the DVM module to validate the ticket internals and to retrieve the
  contents of the ticket
- If the ticket is valid, check that the sport-event already exists or not
- If the sport-event already exists, the event status should not be
  ```status != types.SportEventStatus_STATUS_PENDING```
- Then update the sport-event

---

## **Resolve SportEvent**

When this is processed:

For each event:

- Validate the creator address and validate the ticket format
- Call the DVM module to validate the ticket internals and to retrieve the
  contents of the ticket
- If the ticket is valid, check that the event already exist or not
- If the event already exists, the event status should be
  ```status == types.SportEventStatus_STATUS_PENDING```
- Then resolve the sport-event

```go
resolvedEvent := types.ResolutionEvent{
 Uid            : <string>
 ResolutionTs   : <uint64>
 WinnerOddsUIDs : <map>[string][]byte
 Status         : <SportEventStatus>
}
```

---
