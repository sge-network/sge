# **Accounts**

There is one account in the Bet module.

- Betting Fee Collector: This account holds the betting fee transferred from the bettor to the `bet_fee_collector` module account.
- Price Lock Fee Collector: `price_lock_fee_collector` Temporarily holds the deducted fees of a bet with price-lock feature enabled, the balance will be transferred to the market creator after settlement.
- Price Lock Pool: `price_lock_pool` Temporarily holds the locked balance of the price lock functionality, after settlement, the added amount will be transferred to the bettor if the bet is a winner.

During bet placement, the betting fee is transferred from the bettor's account to the bet module account in the Bet module.

## Wager Transfer

The bet amount (deducted from the betting fee) is transferred to the `orderbook_liquidity_pool` module account in `orderbook` module.
The price lock fee will be deducted from the bettor's wallet balance, this will be temporarily transferred to the `price_lock_fee_collector`.

## Settlement Transfer

- If the user is a winner, the bet amount and payout profit will be transferred from `orderbook_liquidity_pool` module account to the winner's account.
- If the user is a loser, the bet amount and fee will not go back to the bettor's account.
- If the user is eligible to receive price reimbursement, the calculated amount will be transferred to the bettor's account.
- The maximum possible amount of price reimbursement according to the highest price fluctuation will be deducted from the market creator's balance, then in the `orderbook` settlement, the remaining amount will be returned to the market creator's account.
- The deducted fees of price lock will be transferred to the market creator's balance.
