# **Concepts**

The House is tasked with accepting deposits from the users and faciliatates them to act as a house
for a specific market for a small amount of fees. A portion of these fees will be consumed to cover
the transaction costs associated with posting/resolving markets.

The user can withdraw the deposited tokens at any point of time. The tokens that has been
used to accept bets for the house, along with the unused tokens if not withdrawn, will be settled at the
resolution of the market.

## **KYC Validation**

- If Ignore is false in deposit/withdraw ticket payload, then the status of kyc approval should be true and tx signer and kyc id should be same for a deposit/withdraw to be set.
- If Ignore is true in deposit/withdraw ticket payload, then kyc validation is not required and deposit/withdraw can be happen without kyc check.
