# **Concepts**

The **Bet module** is responsible for handling the placement and settlement of bets. Users can place bets on the same market multiple times, either with different or equal odds, using the command line or by singing and broadcasting the bet placement message.

Here are some key points about the betting process:

1. **Wagering and Ticket Verification**:
   - Wagering is facilitated through a ticket that contains odds information, signed by a trusted source.
   - The verification of this ticket occurs during the placement state itself.

2. **Minimum Bet Amount**:
   - Each market defines a minimum bet amount.
   - Users cannot place bets below this minimum threshold.
   - A module parameter controls this requirement.

3. **Betting Fee**:
   - For each market, a specific betting fee has been established.
   - Again, a module parameter governs this fee.

---

**Before Accepting a Bet, Validation Steps:**

1. **Market-Level Validation:**
   - Ensure the market is active for bet placement.
   - Verify that the market is neither resolved nor canceled.
   - Confirm that the bet amount (after deducting the betting fee) meets the minimum allowed bet amount.

2. **Bet-Level Validation:**
   - Validate the provided UUID.
   - Check odds values and apply validations based on American, British, and decimal odds.
   - Ensure the max loss multiplier is positive and less than 1.

3. **OVM-Level Validation:**
   - Validate all data provided in the placement request, including odds values.

4. **KYC Validation:**
   - If the "Ignore" flag in the bet ticket payload is false:
     - Verify that KYC approval status is true.
     - Confirm that the transaction signer and KYC ID match for the bet to be placed.
   - If the "Ignore" flag is true, KYC validation is not required, and the bet can proceed without further checks.

**Wager Assumptions:**

- Users can request to place a single bet only.
- The creator of the transaction is considered the owner of the bet.

**After Bet Acceptance:**

- The bet amount is transferred to the `orderbook_liquidity_pool` module account by the `orderbook` module.
- Betting fees are transferred to the `bet_fee_collector` module account, also managed by the `orderbook` module.
- Bet fulfillments are processed by the `orderbook` module using the `ProcessWager` keeper's method.

## Supported Odds Types

**Note:** Let `bet_amount` be 3,564,819.

- **Decimal (European):** Calculated as `bet_amount * oddsValue`. For example, `3,564,819 * 1.29 = 4,598,616.51`.
- **Fractional (British):** Calculated as `bet_amount + (bet_amount * fraction)`. For instance, `3,564,819 + (3,564,819 * 2/7) = 4,583,338.71`.
- **Moneyline (American):** Calculated as follows:
  - For positive odds value: `bet_amount + (bet_amount * |oddsValue/100|)`. For instance, with odds of +350, `3,564,819 + 3,564,819 * |350/100| = 16,041,685.50` (rounded down).
  - For negative odds value: `bet_amount + (bet_amount * |100/oddsValue|)`. For example, with odds of -350, `3,564,819 + 3,564,819 * |100/-350| = 4,583,338.71` (rounded down).

### Precision

Some online calculators employ a two-digit precision when rounding division results in Fractional and Moneyline calculations. Specifically, these tools convert Moneyline and Fractional odds to Decimal odds and then compute the payout based on the rounded decimal value. This approach significantly impacts the resulting payout.

SGE-Network currently accepts bets with the `usge` currency, which may have substantial value in the market. For such high-value scenarios, it is advisable to perform high-precision calculations within the blockchain code.

> **Note:** The final calculated payout amounts are rounded to two-digit float values, resulting in a slight loss of benefits or payouts.
