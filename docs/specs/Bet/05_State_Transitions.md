# **State Transitions**

This section defines the state transitions of the bet module's KVStore in all scenarios:

## **Place bet**

When this is processed:

- If the ticket is valid a new bet will be created with the given data and will be added to the `bet module's KVStore`.
- The pending bet, ID map and statistics will update accordingly.
- Strategic Reserve module bet placement processor will calculate and transfer bet amount to the corresponding module account.

```go
newBet := &types.Bet{
    Creator:            msg.Creator,
    UID:                msg.UID,
    MarketUID:          <msg.Ticket.MarketUID>,
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
    BetFulfillment:     <bet fulfilment by the strategic reserve>
}
```

---

## **Settle bet**

When this  is processed:

- If corresponding market is aborted or canceled, the bet will be updated in the `bet module's KVStore` as below:

    ```go
    bet.Result = types.Bet_RESULT_REFUNDED
    bet.Status = types.Bet_STATUS_SETTLED
    ```

- Resolve the bet result based on the market result, and update field `Result` to indicate won or lost, and field `Status` to indicate result is declared. For Example:

    ```go
    bet.Result = types.Bet_RESULT_WON // or types.Bet_RESULT_LOST
    bet.Status = types.Bet_STATUS_RESULT_DECLARED
    ```

- Call `Strategic Reserve module` to unlock fund and payout user based on the bet's result, and update the bet's `Status` field to indicate it is settled:

    ```go
    bet.Status = types.Bet_STATUS_SETTLED
    ```

- Store the updated bet in the `bet module's KVStore`.

---

## **Batch bet settlement**

Batch bet settlement happens in the end-blocker of the bet module:

1. Get resolved markets that have the unsettled bets.
    - for each market:
        1. Settle the bets one by by querying the market bets.
        2. Remove the resolved market from the list if there is no more active bet.
        3. Call strategic reserve's method to set the order book as settled
2. Check the `BatchSettlementCount` parameter of bet module and let the rest of bets for the nex block.
