# **Concepts**

The **Subaccount Module** is responsible for creating, topping up, and managing the funds of sub accounts.

## Sub Account

A sub account is a special type of account whose balance can be manipulated by the blockchain code logic. Unlike regular accounts, it cannot be accessed or transferred using the **Bank** or **Staking** modules. Each sub account has a one-to-one relationship with its owner account, meaning a normal account can either be associated with a sub account or have no sub account at all.

## TopUp

Sub Account balances can be topped up on demand using the transaction endpoints provided by the sub account module.

## Withdrawal

Unlocked balances from a sub account can be withdrawn to the owner's main account balance. This operation is facilitated through the transaction endpoints of the sub account module.

## Wager

Owners can use the **Subaccount Module** transaction message interfaces to wager using the sub account module's balance. This process automatically withdraws both unlocked and locked balances from the sub account, utilizing the main account address and balance for placing bets.

## House Deposit/Withdraw

The **Subaccount Module** also allows owners to deposit to or withdraw from the house module using the sub account module's balance.
