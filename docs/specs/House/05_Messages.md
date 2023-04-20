# **Messages**

In this section, we describe the processing of the House messages.

## **MsgDeposit**

Within this message, the user specifies the deposit information they wish to make.

```proto
// Msg defines the house Msg service.
service Msg {
  // Deposit defines a method for performing a deposit of coins to become part of the house corresponding to a market.
  rpc Deposit(MsgDeposit) returns (MsgDepositResponse);
}
```

```proto
// MsgDeposit defines a SDK message for performing a deposit of coins to become
// part of the house corresponding to a market.
message MsgDeposit {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;

  // creator is the account who makes a deposit
  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];
  // market_uid is the uid of market/order book against which deposit is being
  // made.
  string market_uid = 2 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];
  // amount is the amount being deposited on an order book to be a house
  string amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // ticket is the jwt ticket data.
  string ticket = 4;
}

// MsgDepositResponse defines the Msg/Deposit response type.
message MsgDepositResponse {
  // market_uid is the uid of market/order book against which deposit is being
  // made.
  string market_uid = 1 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];
  // participation_index is the index corresponding to the order book
  // participation
  uint64 participation_index = 2
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
}
```

## **Deposit Ticket Payload**

This ticket is being used for validating the KYC.

```proto
// DepositTicketPayload indicates data of the deposit ticket.
message DepositTicketPayload {
  // kyc_data contains the details of user kyc.
  sgenetwork.sge.type.KycDataPayload kyc_data = 1
      [ (gogoproto.nullable) = false ];
}
```


### **Deposit Failure cases**

The transaction will fail if:

- Basic validation fails:
  - Invalid creator address
  - Empty or invalid market uid
  - Invalid amount

---

## **MsgWithddraw**

Within this message, the user provides a deposit UID they wish to make a withdrawal against.

```proto
// Msg defines the Msg service.
service Msg {
  // Withdraw defines a method for performing a withdrawal of coins against a deposit.
  rpc Withdraw(MsgWithdraw) returns (MsgWithdrawResponse);
}
```

```proto
// MsgWithdraw defines a SDK message for performing a withdrawal of coins of
// unused amount corresponding to a deposit.
message MsgWithdraw {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;

  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];
  string market_uid = 2 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];
  // participation_index is the index corresponding to the order book
  // participation
  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
  // mode is the withdrawal mode. It can be full or partial withdraw
  WithdrawalMode mode = 4 [ (gogoproto.moretags) = "yaml:\"mode\"" ];
  // amount is the requested withdrawal amount
  string amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  // ticket is the jwt ticket data.
  string ticket = 6;
}

// MsgWithdrawResponse defines the Msg/Withdraw response type.
message MsgWithdrawResponse {
  // id is the unique identifier for the withdrawal
  uint64 id = 1 [
    (gogoproto.customname) = "ID",
    (gogoproto.jsontag) = "id",
    json_name = "id",
    (gogoproto.moretags) = "yaml:\"id\""
  ];
  // market_uid is the id of market/order book from which withdrawal is made
  string market_uid = 2 [
    (gogoproto.customname) = "MarketUID",
    (gogoproto.jsontag) = "market_uid",
    json_name = "market_uid"
  ];
  // participation_index is the index in order book from which withdrawal is
  // made
  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
}
```

## **Withdraw Ticket Payload**

This ticket is being used for validating the KYC.

```proto
// WithdrawTicketPayload indicates data of the withdrawal ticket.
message WithdrawTicketPayload {
  // kyc_data contains the details of user kyc.
  sgenetwork.sge.type.KycDataPayload kyc_data = 1
      [ (gogoproto.nullable) = false ];
}
```
