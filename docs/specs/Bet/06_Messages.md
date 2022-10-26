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

  // uid is the unique uuid assigned to bet
  string uid = 2;

  // amount is the wagger amount
  string amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  // ticket is a signed string containing important info such as `oddsValue`
  string ticket = 4;
}

// MsgPlaceBetResponse is the returning value in the response of MsgPlaceBet request
message MsgPlaceBetResponse {}
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

## **MsgPlaceBetSlip**

Within this message, the user specifies the bet slip information they wish to place.

```
// PlaceBetSlip defines a method to place multiple bets with the given data
  rpc PlaceBetSlip(MsgPlaceBetSlip) returns (MsgPlaceBetSlipResponse);
```

```
message MsgPlaceBetSlip {
  // creator is the bettor address
  string creator = 1;

  // bets is an array of bet to place
  repeated Bet bets = 2;
}

// MsgPlaceBetSlipResponse is the returning value in the response of MsgPlaceBetSlip request
message MsgPlaceBetSlipResponse {
  // successful_bet_uids_list is an array of successful bet UIDs to place
  repeated string successful_bet_uids_list = 1;

  // failed_bet_uids_error_map is an map of failed bet UIDs to place alongside their failure messages
  map<string, string> failed_bet_uids_error_map = 2;
}
```

### **Failure cases**
The transaction will fail if:
- Creator address is invalid
- No any bet informatoin is provided to place
- The number of given bets is more that a predefined threshold (This threshold is defined as a constant named `BetPlacementThreshold` in code)

A success `PlaceBetSlip` transaction will return two lists for success and failed bets placement. Each bet will be in the failed list if:
- Validation fails:
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
For each bet:
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

## **MsgSettleBetBulk**

Within this message, the user provides a bet UID they wish to settle its corresponding bet.

```
// SettleBetBulk defines a method to settle multiple given bets
  rpc SettleBetBulk(MsgSettleBetBulk) returns (MsgSettleBetBulkResponse);
```

```

// MsgSettleBetBulk defines a message to settle multiple given bets
message MsgSettleBetBulk {
  // creator is the bettor address
  string creator = 1;

  // bet_uids is an array of uuid of the bets to settle
  repeated string bet_uids = 2;
}

// MsgSettleBetBulkResponse is the returning value in the response of MsgSettleBetBulk request
message MsgSettleBetBulkResponse {
  // successful_bet_uids_list is an array of successful bet UIDs to place
  repeated string successful_bet_uids_list = 1;

  // failed_bet_uids_error_map is an map of failed bet UIDs to place alongside their failure messages
  map<string, string> failed_bet_uids_error_map = 2;
}

```

### **Failure cases**
The transaction will fail if:
- Creator address is Invalid
- No any bet UIDs are provided to settle
- The number of given bet UIDs is more that a predefined threshold (This threshold is defined as a constant named `SettlementUidsThreshold` in code)

A success `SettleBetBulk` transaction will return two lists for success and failed bets settlement. Each bet will be in the failed list if:
- Bet UID in invalid
- There is no matching bet
- Bet is canceled
- Bet is already settled
- Corresponding sport event not found
- Result of corresponding sport event is not declared
- There is an error in SR module functions

### **What Happens**
For each bet:
- If corresponding sport event is aborted or canceled, the bet will be updated in the `bet module's KVStore` as below:
    ```
    bet.Result = types.Bet_RESULT_ABORTED
    bet.Status = types.Bet_STATUS_SETTLED
    ```
- Resolve the bet result based on the sport event result, and update field `Result` to indicate won or lost, and field `Status` to indicate result is declared.
- Call `Strategic Reserve module` to unlock fund and payout user based on the bet's result, and update the bet's `Status` field to indicate it is settled.
- Store the updated bet in the `bet module's KVStore`.
