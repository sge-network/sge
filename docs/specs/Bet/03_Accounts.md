# **Accounts**

There is one account in the Bet module.

- Betting Fee Collector: This account holds the betting fee transferred from the bettor to the `bet_fee_collector` module account.
- Price Lock Pool: `price_lock_pool` contains the balance of the price lock pool, which is allocated for compensating the price variation.

During bet placement, the betting fee is transferred from the bettor's account to the bet module account in the Bet module.

## Wager Transfer

The bet amount (deducted from the betting fee) is transferred to the `orderbook_liquidity_pool` module account in `orderbook` module.

## Settlement Transfer

- If the user is a winner, the bet amount and payout profit will be transferred from `orderbook_liquidity_pool` module account to the winner's account.
- If the user is a loser, the bet amount and fee will not go back to the bettor's account.
- If the user is eligible to receive price reimbursement, the calculated amount will be transferred to the bettor's account.
