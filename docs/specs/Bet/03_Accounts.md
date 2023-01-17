# **Accounts**

There is one account in the Bet module.

- Betting Fee Account: This account holds the betting fee transferred from the bettor to the `bet` module account.

During bet placement, betting fee is transferred from the bettor's account to the bet module account in the Bet module. Then it will be used by the Incentive Pendulum (to be implemented in later phases).

> Note: The bet amount (deducted by betting fee) is transferred to the `bet_reserve` module account in SR module when the bet is placed, and after settlement if the user is winner, the profit will be transferred from `bet_reserve` module account to the winner's account and the payout will be transferred from `bet_reserve` to `sr_pool` module account.
