# **State Transitions**

This section defines the state transitions of the bet module's KVStore in all scenarios:

## **Place bet**

When this is processed:

- If the ticket is valid a new bet will be created with the given data and will be added to the `bet module's KVStore`. Like this:

```go
newBet := &types.Bet{
    Creator:       msg.Creator,
    Uid:           msg.Uid,
    SportEventUid: <msg.Ticket.SportEventUid>,
    OddsUid:       <msg.Ticket.OddsUid>,
    OddsType:      <msg.OddsType>,
    OddsValue:     <msg.Ticket.OddsValue>,
    Amount:        msg.Amount,
    BetFee         <will be calculated>,
    Ticket:        msg.Ticket,
    Status:        types.Bet_STATUS_PLACED
    Result:        types.Bet_RESULT_PENDING
    Verified:      true,
    CreatedAt:     <current timestamp of block time>,
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

- Call `Strategic Reserve module` to unlock fund and payout user based on the bet's result, and update the bet's `Status` field to indicate it is settled:

    ```go
    bet.Status = types.Bet_STATUS_SETTLED
    ```

- Store the updated bet in the `bet module's KVStore`.

---
