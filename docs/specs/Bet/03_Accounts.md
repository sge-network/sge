# **Account Details**

Within the Bet module, there exists a single account:

- **Betting Fee Collector**: This account receives the betting fee transferred from the bettor to the `bet_fee_collector` module account.

When placing a bet, the betting fee is moved from the bettor's account to the bet module account in the Bet module.

## Wager Transfer

The bet amount (after deducting the betting fee) is then transferred to the `orderbook_liquidity_pool` module account within the `orderbook` module.

## Settlement Transfer

- If the user is a winner, both the bet amount and the payout profit are transferred from the `orderbook_liquidity_pool` module account to the winner's account.
- However, if the user is a loser, the bet amount and fee remain with the bettor and do not return to their account.
