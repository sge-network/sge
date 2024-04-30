# **Concepts**

The House serves as a platform for accepting user deposits and acts as a facilitator for specific markets, charging nominal fees. A portion of these fees covers transaction costs related to market posting and resolution.

Users can withdraw their deposited tokens at any time. Tokens used for accepting bets on behalf of the house, along with any unused tokens, are settled upon market resolution.

## **KYC Validation**

- When `Ignore` is set to false in the deposit/withdraw ticket payload, KYC approval status must be true, and the transaction signer's KYC ID should match the deposit/withdrawal request.
- If `Ignore` is true, KYC validation is not required, allowing deposit/withdrawal without further checks.

## **Authorization**

Authorization is verified for each deposit and withdrawal request that includes the `depositor_address` in the ticket. This feature is useful when an account wishes to deposit/withdraw on behalf of other accounts.
To grant deposit/withdrawal permissions, the granter utilizes the `authz` module in the Cosmos-SDK.
After each successful deposit or withdrawal, the spend or withdrawal limit is updated. If the spend limit reaches zero, the grant record is completely removed from the authz state.
