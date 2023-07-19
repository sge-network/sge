# **Overview**

The Bet module is responsible for receiving and processing requests to wager and settle bets. In the case of wagering, it validates the request and places the bet.

For the settlement, blockchain automatically queries resolved markets then for each of these markets, checks the result of the market, determines the bet result, and settles the bet using `orderbook` module.
