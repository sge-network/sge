# **State Transitions**

This section defines the state transitions of the sport module's KVStore in all scenarios:

## **Batch Create SportEvent**

When this is processed:

- Validate the creator address and If the ticket is also in valid format
- Call DVM module to validate the ticket internals and to retrieve the content
- If the ticket is valid, iterate over every single entry
- check if event should not already exist
- After all the above steps a new sport event will be created with the given data and will be added to the `sportevent module's KVStore`. Like this:

```
newEvent := &type.SportEvent{
	Uid            : <string>
	StartTS        : <uint64>
	EndTS          : <uint64>
	OddsUIDs       : <[]string>
	WinnerOddsUIDs : <map[string][]byte>
	Status         : <SportEventStatus>
	ResolutionTS   : <uint64>
	Creator        : <string>
	BetConstraints : <*EventBetConstraints>
	Active         : <bool>
}
```
---

## **Batch Update SportEvent**

When this is processed:

For each event:
- Validate the creator address and If the ticket is also in valid format
- Call DVM module to validate the ticket internals and to retrieve the content
- If the ticket is valid, iterate over every single entry
- Check that event should already be created
- Event status cannot be
  ```status != types.SportEventStatus_STATUS_PENDING```
- Afterwards update the information in the existing call

```
k.Keeper.UpdateSportEvent(ctx, types.SportEvent{
			Uid:            sportEvent.Uid,
			StartTS:        event.StartTS,
			EndTS:          event.EndTS,
			BetConstraints: &types.EventBetConstraints{
				MaxBetCap:          event.BetConstraints.MaxBetCap,
				MinAmount:          event.BetConstraints.MinAmount,
				BetFeePercentage:   event.BetConstraints.BetFeePercentage,
			},
			Active: event.Active,
		})
```

---

## **Batch Resolve SportEvent**

When this is processed:

For each event:
- Validate the creator address and If the ticket is also in valid format
- Call DVM module to validate the ticket internals and to retrieve the content
- If the ticket is valid, iterate over every single entry
- Check that event should already be created
- Event status should be
  ```status == types.SportEventStatus_STATUS_PENDING```
- Afterwards resolve the information in the existing call

```
resolvedEvent := types.ResolutionEvent{
	Uid            : <string>
	ResolutionTs   : <uint64>
	WinnerOddsUIDs : <map>[string][]byte
	Status         : <SportEventStatus>
}
```

---
