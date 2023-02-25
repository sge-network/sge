# **State Transitions**

This section defines the state transitions of the bet module's KVStore in all scenarios:

## **Place bet**

When this is processed:

- If the ticket is valid a new bet will be created with the given data and will be added to the `bet module's KVStore`.
- Order Book module bet placement processor will calculate and transfer bet amount to the corresponding module account.

```go
newBet := &types.Bet{
    Creator:            msg.Creator,
    UID:                msg.UID,
    SportEventUID:      <msg.Ticket.SportEventUID>,
    OddsUID:            <msg.Ticket.OddsUID>,
    OddsType:           <msg.OddsType>,
    OddsValue:          <msg.Ticket.OddsValue>,
    Amount:             msg.Amount,
    BetFee:             <will be calculated>,
    Ticket:             msg.Ticket,
    Status:             types.Bet_STATUS_PLACED
    Result:             types.Bet_RESULT_PENDING
    Verified:           true,
    CreatedAt:          <current timestamp of block time>,
    MaxLossMultiplier:  <the coefficient of multiplicitation of the maximum loss>,
    BetFulfillment:     <bet fulfilment by the order book>
}
```

---

## **Settle bet**

When this  is processed:

- If corresponding sport-event is aborted or canceled, the bet will be updated in the `bet module's KVStore` as below:

    ```go
    bet.Result = types.Bet_RESULT_ABORTED
    bet.Status = types.Bet_STATUS_SETTLED
    ```

- Resolve the bet result based on the sport-event result, and update field `Result` to indicate won or lost, and field `Status` to indicate result is declared. For Example:

    ```go
    bet.Result = types.Bet_RESULT_WON
    bet.Status = types.Bet_STATUS_RESULT_DECLARED
    ```

- Call `Order Book module` to unlock fund and payout user based on the bet's result, and update the bet's `Status` field to indicate it is settled:

    ```go
    bet.Status = types.Bet_STATUS_SETTLED
    ```

- Store the updated bet in the `bet module's KVStore`.

---

## **Batch bet settlement**

Batch bet settlement happens in the end-blocker of the bet module:

1. Get resolved sport events that have the unsettled bets.
    - for each sport-event:
        1. Settle the bets one by by querying the sport-event bets.
        2. Remove the resolved sport event from the list if there is no more active bet.
        3. Call order book method to set the order book as settled
2. Check the `BatchSettlementCount` parameter of bet module and let the rest of bets for the nex block.
