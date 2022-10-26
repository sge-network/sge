# **State Transitions**

This section defines the state transitions of the Params

## **Lock funds of the SR pool (called by the `bet` module)**

When a bet is placed by a user, the potential winning payout amount is locked in the SR pool and is made unavailable for accepting further bets.

When this is processed:

- Transfer funds from the user’s account to the SR module account

- update the `locked_amount` and `unlocked_amount` params

---

## **Payout a winning user (called by the `bet` module)**

If the House loses the wager, the payout is made from the Strategic Reserve to the winner’s account directly.

When this  is processed:

- Unlock the corresponding locked amount in the SR

- Transfer funds from the SR pool to the user

- Modify `locked_amount` and `unlocked_amount` params

---

## **Unlock SR pool funds (called by the `bet` module)**

If the House wins a wager, the funds locked in the Strategic Reserve for the corresponding bet are unlocked and transferred to the module account.

When this  is processed:

- Unlock the corresponding locked amount in the SR

- Transfer the amount to the SR module account

- Modify `locked_amount` and `unlocked_amount` params

---
