# **Accounts**

Within the **Reward Module**, there exists a solitary account.

- Reward Pool: This particular account stores the sum of rewards transferred from the `promoter` campaign account to the `reward_pool` module account.

Upon campaign creation, the pool sum will be directed to the reward pool module account.

## **Transfer of Granted Rewards**

The rewards are then transferred to the main (or sub) account of the reward receiver, based on the ticket payload of the grant reward transaction endpoint.
