# **State Transitions**

## **Wager State Transitions**

When processing the wager, the following state transitions occur in the `Bet` module's KVStore:

1. If the ticket is valid, a new bet is created using the provided data and added to the `Bet` module state.
2. The pending bet, ID map, and statistics are updated accordingly.
3. The `orderbook` module's bet placement processor calculates and transfers the bet amount and fee to the corresponding module accounts.

```go
newBet := &types.Bet{
    Creator:            msg.Creator,
    UID:                msg.UID,
    MarketUID:          <msg.Ticket.MarketUID>,
    OddsUID:            <msg.Ticket.OddsUID>,
    OddsValue:          <msg.Ticket.OddsValue>,
    Amount:             msg.Amount,
    BetFee:             <will be calculated>,
    Status:             types.Bet_STATUS_PLACED
    Result:             types.Bet_RESULT_PENDING
    CreatedAt:          <current timestamp of block time>,
    MaxLossMultiplier:  <the coefficient of multiplicitation of the maximum loss>,
    BetFulfillment:     <bet fulfilments by the orderbook>
}
```

---

## **Bet Settlement Process**

When processing a bet:

1. If the corresponding market is aborted or canceled, update the `Bet` module state as follows:

    ```go
    bet.Result = types.Bet_RESULT_REFUNDED
    bet.Status = types.Bet_STATUS_SETTLED
    ```

    Additionally, invoke the `RefundBettor` method from the `orderbook` module to refund the bet amount and bet fee.

2. Determine the bet result based on the market outcome. Update the `Result` field to indicate whether the bet was won or lost, and set the `Status` field to reflect that the result has been declared. For example:

    ```go
    bet.Result = types.Bet_RESULT_WON // or types.Bet_RESULT_LOST
    bet.Status = types.Bet_STATUS_RESULT_DECLARED
    ```

3. Call the `BettorLoses` or `BettorWins` methods from the `orderbook` module to unlock funds and pay out users based on the bet result. Update the bet's `Status` field to indicate that it has been settled:

    ```go
    bet.Status = types.Bet_STATUS_SETTLED
    ```

4. If the market result is declared, invoke the `WithdrawBetFee` method from the `orderbook` module to transfer the bet fee to the market creator's account balance.

5. Finally, store the updated bet information in the state.

---

## Batch Bet Settlement

Batch bet settlement occurs during the final block of the bet module. Follow these steps:

1. Identify resolved markets that still have unsettled bets.
    - For each market:
        1. Settle the bets by querying the market's bet data.
        2. If there are no more active bets in the market, remove it from the list.
        3. Use the orderbook's win/lose methods to transfer the appropriate amounts.
2. Verify the `BatchSettlementCount` parameter in the bet module and handle the remaining bets in the next block.
