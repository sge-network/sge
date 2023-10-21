# **State Transitions**

This section defines the state transitions of the Sub Account module's KVStore in all scenarios:

## **Create Sub Account**

When this is processed:

- Sum the unlock balances according to the input locked balances.
- Check if there is an existing subaccount for the proposed owner address.
- Generate new account using account module of the Cosmos-SDK.
- Transfer the calculated balance to the newly created account.
- Set Sub account owner in the state.
- Set locked balance in the state.
- Set balance in the state.

---

## **Top Up Sub Account**

When this is processed:

- Sum the unlock balances according to the input locked balances.
- Get subaccount by owner address from the state.
- Increase the deposited amount of the balance.
- Set locked balance in the state.
- Set balance in the state.

---

## **Wager**

When this is processed:

- Get Subaccount by Owner Address from the state.
- Call bet module's method to prepare the bet object.
- Get subaccount balance from the state.
- Deduct the bet amount from the sub account balance.
- Call bet module's wager method to set the bet.
- Set the new balance of the sub account module in the state.

---

## **House Deposit**

When this is processed:

- Get Subaccount by Owner Address from the state.
- Get subaccount balance from the state.
- Call house module's method to parse the ticket and valudate.
- Deduct the deposit amount from the sub account balance.
- Call house module's Deposit method to set the participation.
- Set the new balance of the sub account module in the state.

## **House Withdraw**

When this is processed:

- Get Subaccount by Owner Address from the state.
- Get subaccount balance from the state.
- Call house module's method to parse the ticket and valudate.
- Return the withdrawal amount from the sub account balance.
- Call house module's Deposit method to set the participation.
- Set the new balance of the sub account module in the state.
