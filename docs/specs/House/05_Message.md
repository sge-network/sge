# **Messages**

In this section, we describe the processing of the House messages.

## **MsgDeposit**

Within this message, the user specifies the deposit information they wish to make.

```proto
// Msg defines the house Msg service.
service Msg {
  // Deposit defines a method for performing a deposit of coins to become part of the house corresponding to a sport event.
  rpc Deposit(MsgDeposit) returns (MsgDepositResponse);
}
```

```proto
// MsgDepositRequest defines a SDK message for performing a deposit of coins to become
// part of the house corresponding to a sport event.
message MsgDeposit {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;

  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];
  string sport_event_uid = 2 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];
  string amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}

// MsgDepositResponse defines the Msg/Deposit response type.
message MsgDepositResponse {
  string sport_event_uid = 1 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];

  uint64 participation_index = 2
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
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
// MsgWithdraw defines a SDK message for performing a withdrawal of coins of
// unused amount corresponding to a deposit.
message MsgWithdraw {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;

  string creator = 1 [ (gogoproto.moretags) = "yaml:\"creator\"" ];
  string sport_event_uid = 2 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];
  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
  WithdrawalMode mode = 4 [ (gogoproto.moretags) = "yaml:\"mode\"" ];
  string amount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}

// MsgWithdrawResponse defines the Msg/Withdraw response type.
message MsgWithdrawResponse {
  uint64 id = 1 [
    (gogoproto.customname) = "ID",
    (gogoproto.jsontag) = "id",
    json_name = "id",
    (gogoproto.moretags) = "yaml:\"id\""
  ];

  string sport_event_uid = 2 [
    (gogoproto.customname) = "SportEventUID",
    (gogoproto.jsontag) = "sport_event_uid",
    json_name = "sport_event_uid"
  ];

  uint64 participation_index = 3
      [ (gogoproto.moretags) = "yaml:\"participation_index\"" ];
}
```
