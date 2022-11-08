# **Messages**

In this section, we describe the processing of the Bet messages.

## **MsgPlaceBet**

Within this message, the user specifies the bet information they wish to place.

```
// PlaceBet defines a method to place a bet with the given data
  rpc PlaceBet(MsgPlaceBet) returns (MsgPlaceBetResponse);
```

```
// MsgPlaceBet defines a message to place a bet with the given data
message MsgPlaceBet {
  // creator is the bettor address
  string creator = 1;

  // bet is the info of bet to place
  BetPlaceFields bet = 2;
}

// PlaceBetFields contains necessary fields which come in BetPlacement and BetSlipPlacement TX requests
message BetPlaceFields {
  // uid is the unique uuid assigned to bet
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];

  // amount is the wagger amount
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  // ticket is a signed string containing important info such as `oddsValue`
  string ticket = 3;
  
  // odds_type is the type of odds used has chosen
  sgenetwork.sge.bet.OddsType odds_type = 4;
}

// MsgPlaceBetResponse is the returning value in the response of MsgPlaceBet request
message MsgPlaceBetResponse {
    string error = 1;
    BetPlaceFields bet = 2;
}
```

### **Failure cases**
The transaction will fail if:
- Basic validation fails:
  - Invalid creator address
  - Empty or invalid bet UID
  - Empty amount
  - Non positive amount
  - Empty or invalid ticket (containing space)
- Provided bet UID is already set
- Empty or invalid odds UID in ticket
- Empty, non positive or invalid odds value in ticket
- Empty or invalid sport event UID in ticket
- There is no any sport event with the given sportEventUID
- Sport event is not active for accepting bet (it's not active or status in not `PENDING`)
- Sport event has expired
- Sport event maximum betting capacity has been reached
- The sport event does not contain the selected odds
- Bet amount is less than minimum allowed amount
- The creator address is not valid
- There is an error in AddExtraPayoutToEvent in sportEvent module
- There is an error in ProcessBetPlacement in SR module

### **What Happens**
- A new bet will be created with the given data and will be added to the `bet module's KVStore`.
---

## **MsgSettleBet**

Within this message, the user provides a bet UID they wish to settle its corresponding bet.

```
// SettleBet defines a method to settle the given bet
  rpc SettleBet(MsgSettleBet) returns (MsgSettleBetResponse);
```

```
// MsgSettleBet defines a message to settle the given bet
message MsgSettleBet {
  // creator is the bettor address
  string creator = 1;

  // bet_uid is the unique uuid of the bet to settle
  string bet_uid = 2;
}

// MsgSettleBetResponse is the returning value in the response of MsgSettleBet request
message MsgSettleBetResponse {
    string error = 1;
    string bet_uid = 2 [(gogoproto.customname) = "BetUID" ,(gogoproto.jsontag) = "bet_uid", json_name = "bet_uid"];
}
```

### **Failure cases**
The transaction will fail if:
- Basic validation fail:
  - Invalid creator address
  - Empty bet UID
- Bet UID in invalid
- There is no matching bet
- Bet is canceled
- Bet is already settled
- Corresponding sport event not found
- Result of corresponding sport event is not declared
- There is an error in SR module functions

### **What Happens**
- If corresponding sport event is aborted or canceled, the bet will be updated in the `bet module's KVStore` as below:
    ```
    bet.Result = types.Bet_RESULT_ABORTED
    bet.Status = types.Bet_STATUS_SETTLED
    ```
- Resolve the bet result based on the sport event result, and update field `Result` to indicate won or lost, and field `Status` to indicate result is declared.
- Call `Strategic Reserve module` to unlock fund and payout user based on the bet's result, and update the bet's `Status` field to indicate it is settled.
- Store the updated bet in the `bet module's KVStore`.

---
