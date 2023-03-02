# **State Transitions**

This section defines the state transitions of the sport module's KVStore in all scenarios:

## **Add SportEvent**

When this is processed:

- Validate the creator address and validate the ticket format.
- Call the DVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check if the event already exists.
- Use default bet constraints if the `bet_fee` and `min_bet_amount` are not available in the ticket.
- Initiate order book for the sport-event.
- If the sport-event does not already exist, a new sport-event should be
  created with the given data and will be added to the module state.

```go
newEvent := &type.SportEvent{
 Uid                    : <string>
 StartTS                : <uint64>
 EndTS                  : <uint64>
 Odds                   : <[]*Odds>
 WinnerOddsUIDs         : <map[string][]byte>
 Status                 : <SportEventStatus>
 ResolutionTS           : <uint64>
 Creator                : <string>
 BetConstraints         : <*EventBetConstraints>
 Meta                   : <string>
 SrContriButionForHouse : <sdk.Int>
 BookID                 : <string>
}
```

---

## **Update SportEvent**

When this is processed:

For each event:

- Validate the creator address and validate the ticket format.
- Call the DVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check that the sport-event already exists or not.
- The event status should be active or inactive to be updatable, if not
returns appropriate error.
- Then update the sport-event in the module state.

---

## **Resolve SportEvent**

When this is processed:

For each event:

- Validate the creator address and validate the ticket format.
- Call the DVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check that the event already exist or not.
- The sport-event should exist and the status should be active otherwise proper error returned.
- Then resolve the sport-event and set in the module state.

```go
resolvedEvent := types.ResolutionEvent{
 Uid            : <string>
 ResolutionTs   : <uint64>
 WinnerOddsUIDs : <map>[string][]byte
 Status         : <SportEventStatus>
}
```
