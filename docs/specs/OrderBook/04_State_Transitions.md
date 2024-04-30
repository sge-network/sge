# **State Transitions**

This section outlines the state transitions of the **KVStore** in the **Orderbook Module** across various scenarios:

## **Order Book Initialization**

Upon creation of a market:

1. Establishes an order book for the respective market.
2. Sets the exposures for the bet odds.

```go
newOrderBook := &types.Deposit{
    OrderBookID:            <marketID>,
    ParticipantCount:       0, // no participation yet
    OddsCount:              <Count of bet odds of market>
    Status:                 Active
}
```

---

## **Initiate Participation**

Upon depositing tokens:

1. Fetch the market and its order book.
2. Verify market activity; return error if inactive.
3. Confirm order book activity; return error if inactive.
4. Check if the maximum allowed participants have been reached.
5. Set participation equal to the deposited liquidity amount.
6. Transfer the liquidity amount to the `orderbook_liquidity_pool` module account.
7. Transfer the deposition fee to the `house_fee_collector` module account.
8. Update the order book odds exposures and add the participation to the fulfillment queue.
9. Initialize the participation exposure as zero for round one and set it to the state.

---

## **Process Wager**

1. Retrieve order book and odds exposures.
2. Iterate through all fulfillment queue items:
    1. Obtain participations and participation exposures.
    2. Check available liquidity and process fulfillment.
    3. Update the state with participations and exposures.
3. Remove the fulfillment queue item from the order book.
4. Transfer bet fee to the `bet_fee_collector` module account.
5. Transfer fulfilled bet amount to the `orderbook_liquidity_pool` account.
6. Mark the bet as paid out.

## **Process Bet Settlement**

1. When the bettor wins (called by the `bet` module):
    - For each bet fulfillment:
        1. Retrieve participation.
        2. Transfer bet amount and payout profit to the bettor's account address from the `orderbook_liquidity_pool` module account.
        3. Set the actual profit of the participation that fulfilled this bet.
        4. Update the participation in the state.
2. When the bettor loses (called by the `bet` module):
    - For each bet fulfillment:
        1. Retrieve participation.
        2. Update actual profit for the participation.
        3. Update the participation in the state.
3. Refund bettor (called by the `bet` module):
    1. Transfer bet fee to the bettor's account address from the `bet_fee_collector` module account.
    2. Transfer bet amount to the bettor's account address from the `orderbook_liquidity_pool` module account.

---

## **Batch order book settlement**

Batch bet settlement occurs in the end-blocker of the `orderbook` module:

1. Identify resolved order books with no unsettled bets.
    - For each order book (market):
        - If the market is canceled or aborted:
            1. Refund the depositor the original deposit liquidity from the `orderbook_liquidity_pool` module account.
            2. Refund the depositor the original deposit fee from the `house_fee_collector` module account.
            3. Mark the participation as settled in the module state.
        - If the market result is declared and settled:
            1. Refund the depositor the original deposit liquidity plus the actual profit gained in fulfillment from the `orderbook_liquidity_pool` module account.
            2. Refund the depositor the original deposit fee from the `house_fee_collector` module account if the participation did not participate in the bet fulfillment process.
            3. Mark the participation as settled in the module state.
2. Check the `BatchSettlementCount` parameter of the `orderbook` module and proceed with the remaining order books for the next block.
