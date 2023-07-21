# **State Transitions**

This section defines the state transitions of the `market` module state in all scenarios:

## **Add Market**

Validations:

- Validate the creator address and validate the ticket format.
- Call the OVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check if the market already exists.

Modifications:

- If the market does not already exist, a new market should be
  created with the given data and will be added to the module state.
- Initiate order book for the market.

```go
newMarket := &type.Market{
 UID                    : <string>
 StartTS                : <uint64>
 EndTS                  : <uint64>
 Odds                   : <[]*Odds>
 WinnerOddsUIDs         : <map[string][]byte>
 Status                 : <MarketStatus>
 ResolutionTS           : <uint64>
 Creator                : <string>
 Meta                   : <string>
 BookID                 : <string>
}
```

---

## **Update Market**

Validations:

- Validate the creator address and validate the ticket format.
- Call the OVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check that the market already exists or not.
- The market status should be active or inactive to be updatable, if not
returns appropriate error.

Modifications:

- Then update the market in the module state.

---

## **Resolve Market**

Validations:

- Validate the creator address and validate the ticket format.
- Call the OVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check that the market already exist or not.
- The market should exist and the status should be active otherwise proper error returned.

Modifications:

- Then resolve the market and set in the module state.
- Modify list of resolved markets and add newly resolved.

```go
resolvedMarket := types.ResolutionMarket{
 UID            : <string>
 ResolutionTs   : <uint64>
 WinnerOddsUIDs : <map>[string][]byte
 Status         : <MarketStatus>
}
```
