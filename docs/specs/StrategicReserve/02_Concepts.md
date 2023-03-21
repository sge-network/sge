# **Concepts**

The Strategic Reserve is tasked with maintaining the order book, participations, exposures and order book settlement,
each order book will be created as a one-to-one dependency of market. in action, market module calls
book initiation method of the strategic reserve module to create the corresponding order book, participation and exposures.
The created book for the market initiated will be maintained until the market marked as settled.

The Book for a given market, is made up of Book Participants. The first book participant for any book of
a market will be Strategic Reserve. At the creation of a market, in order to facilitate betting on
the created market. Strategic reserve initiates a book for the market and becomes the first participation for the
initiated book.

Once the strategic reserve has initiated a book for a market, users can either bet against the house or
become a part of the house by deposition of chosen amount through the House module. When a user deposits chosen
amount through the house module, the house module will call the strategic reserve module to update the order book
and set the participation for the user on the requested market.

The deposit amount of book participants is used to facilitate betting on the market.

The payout that need to be paid by the system is named Exposure, there are two types of bet odds exposures:

- The odds exposure are the payouts should be expected to be paid.
- The participation exposure are the payout should go to the participant.
