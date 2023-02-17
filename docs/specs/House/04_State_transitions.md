# **State Transitions**

This section defines the state transitions of the house module's KVStore in all scenarios:

## **Despoit**

When this is processed:

- If the sanoity checks for making a deposit passes, a new deposit will be created with the given data and will be added to the `house module's KVStore`. Like this:

```go
newDeposit := &types.Deposit{
    Creator:            msg.Creator,
    Uid:                <will be calculated>,
    SportEventUID:      msg.SportEventUID,
    Amount:             msg.Amount,
    WithdrawnAmount:    nil,
    BetFee              <will be calculated>,
    Status:             types.Deposit_STATUS_ACTIVE,
    CreatedAt:          <current timestamp of block time>,
    SettledAt:          nil,
}
```

---

## **Withdraw**

When this  is processed:

- If the amount is withdrawable, the following changes will be made to the bet-

    ```go
    bet.WithdrawnAmount = bet.WithdrawnAmount + Withdrawal Amount
    ```
