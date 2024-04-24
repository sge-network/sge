# **Overview**

The **Subaccount Module** handles the creation and funding of special accounts called `subaccounts`. These subaccounts have distinct addresses separate from the main (owner) account. Only the **Subaccount module’s** transaction endpoints and events allow for funding or refunding of these accounts. The owner cannot manually transfer funds from this account using the **Bank Module**.

Internally, the Reward module utilizes methods from the **Subaccount Module** to grant rewards to the corresponding subaccount of the reward receiver’s main account.
