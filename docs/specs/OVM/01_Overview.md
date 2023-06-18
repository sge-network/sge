# **Overview**

The SGE-Network chain relies on external off-chain data of matches, markets and KYC validation of the bettors and depositors. To push this data reliably to the chain, some kind of origin verification is required. The `ovm` module essentially fills this role in the SGE-Network chain. The `ovm` module verifies the following data:

- Market data pushed by the House to the chain
- Validity of Odds data using proposed ticket data during bet placement by user
- Depositor's KYC validation
