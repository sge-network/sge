# **State Transitions**

This section defines the state transitions of the Params

## **Process Bet Placement of the SR pool (called by the `bet` module)**

When a bet is placed by a user, the potential winning payout amount is locked in the SR pool and is made unavailable for accepting further bets.

When this is processed:

- Transfer funds from the user’s account to the SR module account

- Update the `locked_amount` and `unlocked_amount` params by locking the
  SR's contribution in the SR Pool

- The bet amount is kept in the `bet_reserve` module account of the SR

- The bet fee is stored in the module account of the bet module

---

## **Bettor wins (called by the `bet` module)**

If the House loses the wager, the payout is made from the Strategic Reserve to the winner’s account directly.

When this  is processed:

- Unlock the corresponding locked amount in the SR

- Transfer funds from the SR pool and the `bet_reserve` to the user

- Modify `locked_amount` and `unlocked_amount` params

---

## **Bettor Loses (called by the `bet` module)**

If the House wins a wager, the funds locked in the Strategic Reserve for the corresponding bet are unlocked and transferred to the module account.

When this  is processed:

- Transfer the bet amount (house winnings) from the `bet_reserve` to the `sr_pool`
  module account of SR

- Unlock the payout profit or SR's contribution in the `sr_pool` by modifying
  the `locked_amount` and `unlocked_amount` params

---

## **Refund Bettor (called by the `bet` module)**

In case the sport-event gets cancelled or aborted, the amount of the bets placed in the system should be returned to the bettor.

When this is processed:

- Refund back the bettor's bet amount from the `bet_reserve` to the bettor in case a sports event gets cancelled or aborted.

---
