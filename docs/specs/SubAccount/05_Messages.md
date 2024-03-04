# **Messages**

In this section, we describe the processing of the Sub Account messages. the transaction message
handler endpoints is as follows

```proto
// Msg defines the Msg service.
service Msg {
  // Create defines a method for creating a subaccount.
  rpc Create(MsgCreate) returns (MsgCreateResponse);

  // TopUp defines a method for topping up a subaccount.
  rpc TopUp(MsgTopUp) returns (MsgTopUpResponse);

  // WithdrawUnlockedBalances defines a method for withdrawing unlocked
  // balances.
  rpc WithdrawUnlockedBalances(MsgWithdrawUnlockedBalances)
      returns (MsgWithdrawUnlockedBalancesResponse);

  // PlaceBet defines a method for placing a bet using a subaccount.
  rpc Wager(MsgWager) returns (MsgWagerResponse);

  // HouseDeposit defines a method for depositing funds to provide liquidity to
  // a market.
  rpc HouseDeposit(MsgHouseDeposit) returns (MsgHouseDepositResponse);

  // HouseWithdraw defines a method for withdrawing funds from a market.
  rpc HouseWithdraw(MsgHouseWithdraw) returns (MsgHouseWithdrawResponse);
}
```

## **MsgCreate**

The sub account creator choose the owner and the initial locked balance.

```proto
// MsgCreate defines the Msg/Create request type.
message MsgCreate {
  // creator is the msg signer.
  string creator = 1;

  // owner is the owner of the subaccount.
  string owner = 2;

  // locked_balances is the list of balance locks.
  // Fixme: why this attribute needs to be repeated?
  repeated LockedBalance locked_balances = 3 [ (gogoproto.nullable) = false ];
}

// MsgCreateAccountResponse defines the Msg/CreateAccount response type.
message MsgCreateResponse {}
```

## **MsgTopUp**

The sub account creator is able to top up the balance of subaccount.

```proto
// MsgTopUp defines the Msg/TopUp request type.
message MsgTopUp {
  // creator is the msg signer.
  string creator = 1;

  // address is the subaccount address.
  string address = 2;

  // locked_balances is the list of balance locks.
  // Fixme: Are we sending multiple balance update together? If not, then only
  // locked balance should be enough
  repeated LockedBalance locked_balances = 3 [ (gogoproto.nullable) = false ];
}

// MsgTopUpResponse defines the Msg/TopUp response type.
message MsgTopUpResponse {}
```

## **MsgWithdrawUnlockedBalances**

The sub account creator is able to withdraw the unlocked balance.

```proto
// MsgWithdrawUnlockedBalances defines the Msg/WithdrawUnlockedBalances request
// type.
message MsgWithdrawUnlockedBalances {
  // creator is the subaccount owner.
  string creator = 1;
}

// MsgWithdrawUnlockedBalancesResponse defines the Msg/WithdrawUnlockedBalances
// response type.
message MsgWithdrawUnlockedBalancesResponse {}
```

## **MsgWager**

The user is able to wager using the sub account balance.

```proto
// MsgWager wraps the MsgWager message. We need it in order not to have
// double interface registration conflicts.
message MsgWager {
  // creator is the subaccount owner.
  string creator = 1;
  // ticket is the jwt ticket data.
  string ticket = 2;
}

// MsgBetResponse wraps the MsgPlaceBetResponse message. We need it in order not
// to have double interface registration conflicts.
message MsgWagerResponse { sgenetwork.sge.bet.MsgWagerResponse response = 1; }
```

## **MsgHouseDeposit**

The user is able to deposit an amount to be a house using the sub account balance.

```proto
// MsgHouseDeposit wraps the MsgHouseDeposit message. We need it in order not to
// have double interface registration conflicts.
message MsgHouseDeposit { sge.house.MsgDeposit msg = 1; }

// MsgHouseDepositResponse wraps the MsgHouseDepositResponse message. We need it
// in order not to have double interface registration conflicts.
message MsgHouseDepositResponse { sge.house.MsgDepositResponse response = 1; }
```

## **MsgHouseWithdraw**

The user is able to withdraw the free deposited (be as a house) amount from the sub account balance.

```proto
// MsgHouseWithdraw wraps the MsgHouseWithdraw message. We need it in order not
// to have double interface registration conflicts.
message MsgHouseWithdraw { sge.house.MsgWithdraw msg = 1; }

// MsgHouseWithdrawResponse wraps the MsgHouseWithdrawResponse message. We need
// it in order not to have double interface registration conflicts.
message MsgHouseWithdrawResponse { sge.house.MsgWithdrawResponse response = 1; }
```
