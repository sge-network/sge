# **Accounts**

There is one account in the Bet module.

- Betting Fee Collector: This account holds the betting fee transferred from the bettor to the `bet_fee_collector` module account.

During bet placement, betting fee is transferred from the bettor's account to the bet module account in the Bet module.

## Wager Transfer

The bet amount (deducted by betting fee) is transferred to the `orderbook_liquidity_pool` module account in `orderbook` module.

## Settlement Transfer

- if the user is winner, the bet amount and payout profit will be transferred from `orderbook_liquidity_pool` module account to the winner's account.
- If the user is loser, the bet amount and fee will not go back to bettor's account.
