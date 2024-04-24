# **Overview**

The SGE-Network chain depends on external off-chain data concerning matches, markets, and KYC validation of both bettors and depositors. Ensuring the reliable transmission of this data to the chain necessitates some form of origin verification. Within the SGE-Network chain, the `ovm` module serves this essential function. It verifies the following types of data:

- Market data transmitted from the House to the chain
- The legitimacy of Odds data by cross-referencing proposed ticket data during user bet placement
- KYC validation of depositors
