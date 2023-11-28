# **Overview**

The Reward module is responsible for Campaign state management and reward validation and application.

To grant rewards, this module uses `subaccount` module methods to TopUp the sub account balance and withdraw it when there is a withdraw request or wager or house deposit.
