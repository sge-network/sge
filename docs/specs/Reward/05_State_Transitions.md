# **State Transitions**

This section defines the state transitions of the reward module's KVStore in all scenarios:

## **Create Campaign**

When this is processed:

- If the ticket is valid a new campaign will be created with the given data and will be added to the `Reward` module state.

```go
newCampaign := &types.Campaign{
    Creator:            msg.Creator,
    UID:                msg.UID,
    Promoter:           msg.ticket.Promoter,
    StartTS:            msg.ticket.StartTs,
    EndTS:              msg.ticket.EndTs,
    RewardCategory:     msg.ticket.RewardCategory,
    RewardType:         msg.ticket.RewardType,
    reward_amount_type: msg.ticket.RewardAmountType,
    reward_amount:      msg.ticket.RewardAmount,
    Pool:               Pool{ Total: msg.Ticket.TotalFunds },
    is_active:          msg.ticket.IsActive,
    claims_per_category:msg.ticket.ClaimsPerCategory,
    meta:               msg.ticket.Meta,
}
```

---

## **Grant Reward**

When this is processed:

- If the corresponding campaign exists and is not expired, continue the process.

- Calculate reward distribution according to the reward amounts defined int the campaign and the reward type and category.  

- Validate availability of the pool balance for the campaign.

- Distribute the rewards according to the calculated distributions.

- Update the campaign pool according to the distribution.

- Set Reward into the module state.

- Set Reward by receiver into the module state.

- Set Reward by campaign into the module state.

> Note: The reward application modifies the campaign pool balance and accounts balances, but does not store reward application in the module state.
