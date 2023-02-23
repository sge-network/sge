# **State Transitions**

This section defines the state transitions of the order book module's KVStore in all scenarios:

## **Book Initiazation**

When this sport-event is being created:

1. Creates a order book corresponding to the sport-event
2. Transfer the sr contibution of sport-event to the liquidity name.
3. Creates the first participation for the sr module account.
4. Set the exposures for the bet odds.
5. Set the exposures for the participations.

```go
newOrderBook := &types.Deposit{
    BookID:                 <sportEventID>,
    ParticipantCount:       1, // this is sr module account
    OddsCount:              <Count of bet odds og sport event>
    Status:                 Active
}
```

---

## **Initiate Participation**

When a user deposits tokens:

1. Retrive the order book.
2. Check if maximum allowed participant is reached or not.
3. Set the participation equal to the liquidiy amount of the deposition.
4. Transfer the deposition fee to the house module account.
5. Transfer the liquidity amount to the Book liquidity module account.
6. Update the book odds exposures and add the participation into the fulfillment queue.
7. Intialize the participation exposure as zero for round as 1 and set to the state.

---

## **Process Bet Placement**

1. Get order book and odds exposures.
2. Check all fulfillment queue items:
    1. Get participationa and participation exposures
    2. Check available liquidity and process fullfilment.
    3. Set the Participation and exposures into the state.
3. Remove the fulfillment queu item from the order book.
4. Transfer bet fee to bet module.
5. Transfer liqudated amount to the bet liquidity module account.
6. Set the bet as paid out bet.

## **Process Bet Settlement**

1. BettorWins(called by the `bet` module):
    - For all bet fulfilments:
        1. Get participation.
        2. Transfer fulfillment payout to the bettor account address from liquidity module account.
        3. Transfer fulfillment bet amount to the bettor account address from liquidity module account.
        4. Set the participation in the state.
2. Bettor Loses(called by the `bet` module):
    - For all bet fulfilments:
        1. Get participation.
        2. Update  actual profit to the paticipations.
        3. Set the participation in the state.

3. Refund bettor(called by the `bet` module):
    - Transfer the bet amount from bet reserve module account to the user account address.
