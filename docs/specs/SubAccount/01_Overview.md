# **Overview**

The Sub Account module is responsible for creating, and funding a special accounts named `subaccount` that have different address from the main(owner) account, this type of account fund/refund is only allowed by this module transaction endpoints and events, the owner is not able to transfer any funds from this account manually using the `bank` module.

Reward module uses this module's methods internally to grant rewards to corresponding sub account of the reward receiver main account.
