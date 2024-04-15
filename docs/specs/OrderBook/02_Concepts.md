# **Concepts**

The **Orderbook Module** is responsible for managing order books, participations, exposures, and order book settlements. Each order book is created as a one-to-one dependency of a market. In practice, the **Market Module** invokes the order book initiation method of the **Orderbook Module** to establish the corresponding order book, participation, and odds exposures. Once created, the order book for an initiated market remains active until the market is marked as settled.

During the order book creation process, there are no available funds to cover the payout profit payment for bets. Therefore, at least one deposit is required on the order book to enable the **Bet Module** to accept bets.

After the **Orderbook Module** initializes an order book for a market, users have two options: they can either bet against the house or become part of the house by depositing a chosen amount through the **House Module**. When a user makes a deposit via the **House Module**, the **Orderbook Module** updates the order book and sets the user's participation for the requested market.

The deposit amounts from order book participants facilitate betting on the market.

System payouts are referred to as "Exposure," which comes in two types of bet odds exposures:

1. **Odds Exposure**: These are the expected payouts.
2. **Participation Exposure**: These payouts are guaranteed based on participation.
