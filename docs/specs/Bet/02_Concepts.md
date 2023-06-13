# **Concepts**

Bet module is tasked with placement and settlement of the bets. the user can place bet on a same market multiple times with different or equal odds through commandline or singing and broadcasting the bet placement message.

> Bet Placement will be done using a ticket containing odds info signed by a trusted source. Verifying this ticket will be done in the placement state itself.
> Bet amount can not be less than a minimum amount which is defined for each market. A module parameter is used for this purpose.
> Also, a betting fee has been defined for each market, A module parameter is used for this purpose.

---

Before accepting bet some validation should take place:

- Market level validation:
  - Market is active for bet placement.
    - Market is not already resolved or canceled.
    - Bet amount (deducted by betting fee) is not less than the minimum allowed bet amount.
- Bet level validation:
  - Provided UUID is valid.
  - Odds Value and the validations according to the American, British and decimal odds.
  - Check if the max loss multiplier is positive and less than 1.
- OVM level validation:
  - All data provided in placement request is valid e.g. odds value.
- KYC Validation:
  - If Ignore is false in bet ticket payload, then the status of kyc approval should be true and tx signer and kyc id should be same for a bet to be placed.
  - If Ignore is true in bet ticket payload, then kyc validation is not required and bet can be placed without kyc check.

Placement Assumptions:

- For bet placement user can raise a request to place a single bet, it can be done for a single bet only.
- When a user is raising a transaction to place a bet, the creator of the transaction is the owner of the  bet.

After a bet is accepted:

- Bet amount transfer to the `orderbook_liquidity_pool` module account this is done by the `orderbook` module.
- Betting fee will be transferred to the `bet_fee_collector` module account. this is done by the `orderbook` module.
- Bet fulfillments are being processed by `orderbook` module in the `ProcessBetPlacement` keeper's method.

## Supported Odds Types

> Note: Let bet_amount be 3564819

- ***Decimal(European):*** Calculated as `bet_amount * oddsValue` ex. `3564819 * 1.29 = 4598616.51`.
- ***Fractional(British):*** Calculated as `bet_amount +  (bet_amount * fraction)` ex. `3564819 + (3564819 * 2/7) = 4583338.71`.
- ***Moneyline(American):*** Calculated as:
  - Positive odds value: `bet_amount + (bet_amount * |oddsValue/100|)` ex. `3564819 + 3564819 * |+350/100| = 16041685.50` the result will be rounded to floor.
  - Negative odds value: `bet_amount + (bet_amount * |100/oddsValue|)` ex. `3564819 + 3564819 * |100/-350| = 4583338.71` the result will be rounded to floor.

### Precision

Some of the Online Calculators round the division result to two-digit precision in Fractional and Moneyline calculations. In other words, these online calculators try to convert Moneyline and Fractional odds to Decimal odds and then calculate the payout according to the calculated rounded decimal value. This approach makes a big difference in the resulting payout. SGE-Network is accepting bets with usge that may have a high value in the market. For this kind of value, it is better to have a high-precision calculation in the blockchain code.

> Note: The final calculated payout amounts are rounded to 2 digit float values, so we have a small portion of lost benefits/payouts.
