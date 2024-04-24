# **State Transitions**

This section outlines the state transitions of the KVStore within the **Subaccount Module** across various scenarios:

## **Creating a Sub Account**

Upon processing this:

- Aggregate the unlocked balances based on the provided locked balances.
- Verify if a `subaccount` already exists for the designated owner address.
- Generate a new account using the Cosmos-SDK account module.
- Transfer the calculated balance to the newly created account.
- Update the state with the sub account owner.
- Update the state with the locked balance.
- Update the state with the balance.

---

## **Topping Up a Sub Account**

Upon processing this:

- Aggregate the unlocked balances based on the provided locked balances.
- Retrieve the `subaccount` associated with the owner address from the state.
- Increment the deposited amount of the balance.
- Update the state with the locked balance.
- Update the state with the balance.

---

## **Wagering**

Upon processing this:

- Retrieve the `subaccount` associated with the owner address from the state.
- Invoke the bet module's method to prepare the bet object.
- Retrieve the balance of the `subaccount` from the state.
- Withdraw locked/unlocked balances of the `subaccount` based on the input proportion.
- Invoke the wager method of the **Bet Module** to establish the bet.
- Update the state with the new balance of the **Subaccount Module**.

---

## **House Deposit**

Upon processing this:

- Retrieve the `subaccount` associated with the owner address from the state.
- Retrieve the balance of the `subaccount` from the state.
- Invoke the method of the **House Module** to parse the ticket and validate.
- Deduct the deposit amount from the sub account balance.
- Invoke the Deposit method of the **House Module** to record the participation.
- Update the state with the new balance of the **Subaccount Module**.

---

## **House Withdrawal**

Upon processing this:

- Retrieve the `subaccount` associated with the owner address from the state.
- Retrieve the balance of the `subaccount` from the state.
- Invoke the method of the **House Module** to parse the ticket and validate.
- Return the withdrawal amount from the sub account balance.
- Invoke the Deposit method of the **House Module** to record the participation.
- Update the state with the new balance of the **Subaccount Module**.
