# **State Transitions**

This section defines the state transitions of the sport module's KVStore in all scenarios:

## **Add Market**

When this is processed:

- Validate the creator address and validate the ticket format.
- Call the OVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check if the market already exists.
- Use default bet constraints if the `bet_fee` and `min_bet_amount` are not available in the ticket.
- Initiate order book for the market.
- If the market does not already exist, a new market should be
  created with the given data and will be added to the module state.

```go
newMarket := &type.Market{
 Uid                    : <string>
 StartTS                : <uint64>
 EndTS                  : <uint64>
 Odds                   : <[]*Odds>
 WinnerOddsUIDs         : <map[string][]byte>
 Status                 : <MarketStatus>
 ResolutionTS           : <uint64>
 Creator                : <string>
 BetConstraints         : <*MarketBetConstraints>
 Meta                   : <string>
 SrContriButionForHouse : <sdk.Int>
 BookID                 : <string>
}
```

---

## **Update Market**

When this is processed:

For each market:

- Validate the creator address and validate the ticket format.
- Call the OVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check that the market already exists or not.
- The market status should be active or inactive to be updatable, if not
returns appropriate error.
- Then update the market in the module state.

---

## **Resolve Market**

When this is processed:

For each market:

- Validate the creator address and validate the ticket format.
- Call the OVM module to validate the ticket internals and to retrieve the
  contents of the ticket.
- If the ticket is valid, check that the market already exist or not.
- The market should exist and the status should be active otherwise proper error returned.
- Then resolve the market and set in the module state.

```go
resolvedMarket := types.ResolutionMarket{
 Uid            : <string>
 ResolutionTs   : <uint64>
 WinnerOddsUIDs : <map>[string][]byte
 Status         : <MarketStatus>
}
```
