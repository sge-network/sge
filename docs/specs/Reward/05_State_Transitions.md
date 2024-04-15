# **State Transitions**

This segment delineates the state transitions of the **KVStore** within the **Reward module** across various situations:

## **Creating a Campaign**

Upon processing:

1. In case of a valid ticket, a new campaign will be generated using the provided data and integrated into the **Reward Module** state.
2. The total pool amount will be subtracted from the promoter's account balance and retained within the pool module account.

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
    meta:               msg.ticket.Meta,
    cap_count:          msg.ticket.CapCount,
    constraints:        msg.ticket.Constraints
}
```

---

## **Grant Reward**

Upon processing:

1. If the corresponding campaign exists and is not expired, proceed with the process.
2. Determine the distribution of rewards based on the defined reward amounts in the campaign, as well as the reward type and category.
3. Verify the availability of the pool balance for the campaign.
4. Allocate the rewards according to the calculated distributions.
5. Update the campaign pool based on the distribution.
6. Record the reward in the module state.
7. Record the reward by receiver in the module state.
8. Record the reward by campaign in the module state.

> Note: The reward application alters the campaign pool balance and accounts balances but does not store the reward application in the module state.
