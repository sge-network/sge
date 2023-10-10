# **State Transitions**

This section defines the state transitions of the reward module's KVStore in all scenarios:

## **Create Campaign**

When this is processed:

- If the ticket is valid a new campaign will be created with the given data and will be added to the `Reward` module state.

```go
newCampaign := &types.Campaign{
    Creator:            msg.Creator,
    FunderAddress:      msg.ticket.FunderAddress,
    UID:                msg.UID,
    StartTS:            msg.ticket.StartTs,
    EndTS:              msg.ticket.EndTs,
    RewardType:         msg.ticket.RewardType,
    RewardDefs:         msg.ticket.RewardDefs,
    Pool:               Pool{ Total: msg.Ticket.PoolAmount },
}
```

---

## **Apply Reward**

When this is processed:

- If the corresponding campaign exists and is not expired, continue the process.

- Calculate reward distribution according to the reward definitions of the campaign and the reward type.  

- Validate availability of the pool balance for the campaign.

- Distribute the rewards according to the calculated distributions.

- Update the campaign pool according to the distribution.

> Note: The reward application modifies the campaign pool balance and accounts balances, but does not store reward application in the module state.
