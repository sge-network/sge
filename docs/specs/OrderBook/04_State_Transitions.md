# **State Transitions**

This section defines the state transitions of the `orderbook` module's KVStore in all scenarios:

## **Order Book Initialization**

When a market is being created:

1. Creates an order book corresponding to the market
2. Set the exposures for the bet odds.

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

When a user deposits tokens:

1. Retrieve the market and the order book.
2. Check if the market is active, if not return error.
3. Check if the order book is active, if not return error.
4. Check if maximum allowed participant is reached or not.
5. Set the participation equal to the liquidity amount of the deposition.
6. Transfer the liquidity amount to the `orderbook_liquidity_pool` module account.
7. Transfer the deposition fee to the `house_fee_collector` module account.
8. Update the order book odds exposures and add the participation into the fulfillment queue.
9. Initialize the participation exposure as zero for round as 1 and set to the state.

---

## **Process Wager**

1. Get order book and odds exposures.
2. Check all fulfillment queue items:
    1. Get participations and participation exposures
    2. Check available liquidity and process fulfillment.
    3. Set the Participation and exposures into the state.
3. Remove the fulfillment queue item from the order book.
4. Transfer bet fee to `bet_fee_collector` module account.
5. Transfer fulfilled bet amount to the `orderbook_liquidity_pool` account.
6. Set the bet as paid out bet.

## **Process Bet Settlement**

1. BettorWins(called by the `bet` module):
    - For all bet fulfilments:
        1. Get participation.
        2. Transfer bet amount and payout profit to the bettor's account address from `orderbook_liquidity_pool` module account.
        3. Set actual profit of the participation that fulfilled this bet.
        4. Set the participation in the state.
2. Bettor Loses(called by the `bet` module):
    - For all bet fulfilments:
        1. Get participation.
        2. Update  actual profit to the paticipations.
        3. Set the participation in the state.

3. Refund bettor(called by the `bet` module):
    1. Transfer bet fee to the bettor's account address from `bet_fee_collector` module account.
    2. Transfer bet amount to the bettor's account address from `orderbook_liquidity_pool` module account.address.

---

## **Batch order book settlement**

Batch bet settlement happens in the end-blocker of the `orderbook` module:

1. Get resolved orderbooks that have no unsettled bets.
    - for each orderbook(market):
        - If the market is canceled or aborted:
            1. Refund depositor the original deposit liquidity from `orderbook_liquidity_pool` module account.
            2. Refund depositor the original deposit fee from `house_fee_collector` module account.
            3. Set the participation as settled in the module state.
        - If market result is declared and settled:
            1. Refund depositor the original deposit liquidity plus the actual profit gained in fulfillment from `orderbook_liquidity_pool` module account.
            2. Refund depositor the original deposit fee from `house_fee_collector` module account if the participation not participated in the bet fulfillment process.
            3. Set the participation as settled in the module state.
2. Check the `BatchSettlementCount` parameter of `orderbook` module and let the rest of order books for the nex block.
