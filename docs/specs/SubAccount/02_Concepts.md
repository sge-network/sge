# **Concepts**

Sub Account module is tasked with create, topup and fund/refund of a sub account.

## **Sub Account**

Sub Account is a special type of account that manipulating its balance can be done by the blockchain code logic and nobody can use `bank` or `staking` modules to transfer or use the subaccount balance to delegate.

> There is a One(None)-To-One relationship between a subaccount and its owner account, so a normal account is able to have a subaccount associated with it or no subaccount associated.

## **TopUp**

Sub Account balance can be topped up on demand and can be done by the transaction endpoints of the sub account module.

## **Withdrawal**

Sub Account unlocked balance can be withdrawn to the owner's account balance. this can be done by the transaction endpoints of the sub account module.

## **Wager**

Using the subaccount module transaction message interfaces, the owner is able to wager using the sub account module's balance.

## **HouseDeposit/Withdraw**

Using the subaccount module transaction message interfaces, the owner is able to deposit/withdraw to the house module using the sub account module's balance.
