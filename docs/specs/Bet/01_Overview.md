# **Overview**

The **Bet module** handles requests related to betting and bet settlement. Let's break it down:

- **Wagering**:
  - When a user places a bet, the Bet module validates the request.
  - It then processes the bet and records it.

- **Settlement**:
  - After a market is resolved, the blockchain automatically queries the results.
  - For each market, the Bet module:
    - Checks the outcome.
    - Determines the bet result.
    - Settles the bet using the `orderbook` module.

In summary, the Bet module ensures smooth betting interactions and accurate settlement based on market outcomes.
