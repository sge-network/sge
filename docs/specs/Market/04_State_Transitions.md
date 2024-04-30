# **State Transitions**

This section outlines the state transitions for the `market` module across various scenarios:

## **Adding a Market**

**Validations:**

- Verify the creator's address and validate the ticket format.
- Utilize the OVM module to validate the ticket internals and retrieve the ticket's contents.
- If the ticket is valid, check whether the market already exists.

**Modifications:**

- If the market does not exist, create a new market with the provided data and add it to the module state.
- Initialize the order book for the market.

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

## **Updating a Market**

**Validations:**

- Validate the creator's address and verify the ticket format.
- Use the OVM module to validate the ticket internals and retrieve the ticket's contents.
- If the ticket is valid, check whether the market already exists.
- Ensure that the market status is either active or inactive for it to be updatable; otherwise, return an appropriate error.

**Modifications:**

- Update the market in the module state.

---

## **Resolving a Market**

**Validations:**

- Validate the creator's address and verify the ticket format.
- Leverage the OVM module to validate the ticket internals and retrieve the ticket's contents.
- If the ticket is valid, check whether the market exists.
- The market must exist, and its status should be active; otherwise, return an appropriate error.

**Modifications:**

- Resolve the market and update it in the module state.
- Adjust the list of resolved markets by adding the newly resolved ones.

```go
resolvedMarket := types.ResolutionMarket{
 UID            : <string>
 ResolutionTs   : <uint64>
 WinnerOddsUIDs : <map>[string][]byte
 Status         : <MarketStatus>
}
```
