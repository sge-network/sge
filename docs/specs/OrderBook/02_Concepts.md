# **Concepts**

The `orderbook` is tasked with maintaining the order book, participations, exposures and order book settlement,
each order book will be created as a one-to-one dependency of market. in action, `market` module calls
order book initiation method of the `orderbook` module to create the corresponding order book, participation and odds exposures.
The created order book for the market initiated will be maintained until the market marked as settled.

When an orderbook is being created, There is no fund available to cover the payout profit payment of the bets, so at least 1 deposit is needed to be made on the order book to enable `bet` module to accept bets.

Once the `orderbook` has initiated an order book for a market, users can either bet against the house or
become a part of the house by deposition of chosen amount through the House module. When a user deposits chosen
amount through the `house` module, the `house` module will call the `orderbook` module to update the order book
and set the participation for the user on the requested market.

The deposit amount of order book participants is used to facilitate betting on the market.

The payout that need to be paid by the system is named Exposure, there are two types of bet odds exposures:

- The odds exposure are the payouts that expected to be paid.
- The participation exposure are the payout that is guaranteed to be paid by the participation.
