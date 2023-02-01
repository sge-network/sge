# **Messages**

In this section, we describe the processing of the House messages.

## **MsgDeposit**

Within this message, the user specifies the deposit information they wish to make.

```proto
// Msg defines the house Msg service.
service Msg {
  // Deposit defines a method for performing a deposit of coins to become part of the house corresponding to a sport event.
  rpc Deposit(MsgDeposit) returns (MsgDepositResponse);

  // Withdraw defines a method for performing a withdrawal of coins against a deposit.
  rpc Withdraw(MsgWithdraw) returns (MsgWithdrawResponse);
}
```

```proto
// MsgDeposit defines a SDK message for performing a deposit of coins to become part of the house corresponding to a sport event.
message MsgDeposit {
  option (gogoproto.equal)           = false;
  option (gogoproto.goproto_getters) = false;
  
  string                   depositor_address = 1 [(gogoproto.moretags) = "yaml:\"depositor_address\""];
  string                   sport_event_uid = 2 [(gogoproto.moretags) = "yaml:\"sport_event_uid\""];
  cosmos.base.v1beta1.Coin amount            = 3 [(gogoproto.nullable) = false];
}

// MsgDepositResponse defines the Msg/Deposit response type.
message MsgDepositResponse {
    string error = 1;
    uint32 deposit_id = 2;
}
```

### **Placement Failure cases**

The transaction will fail if:

- Basic validation fails:
  - Invalid creator address
  - Empty or invalid sport event uid
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
// MsgWithdraw defines a message for performing a withdrawal of coins against a deposit.
message MsgWithdraw {
  // bet_uid is the unique uuid of the bet to settle
  string deposit_uid = 1 [(gogoproto.customname) = "BetUID" ,(gogoproto.jsontag) = "deposit_uid", json_name = "bedeposit_uidt_uid"];

  cosmos.base.v1beta1.Coin amount            = 2 [(gogoproto.nullable) = false];
}

// MsgWithdrawResponse is the returning value in the response of MsgWithdraw request
message MsgWithdrawResponse {
    string error = 1;
    string deposit_uid = 2;
}

```