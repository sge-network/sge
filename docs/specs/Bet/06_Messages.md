# **Messages**

In this section, we describe the processing of the Bet messages. the transaction message 
handler endpoints is as follows

```proto
// Msg defines the Msg service.
service Msg {
  // PlaceBet defines a method to place a bet with the given data
  rpc PlaceBet(MsgPlaceBet) returns (MsgPlaceBetResponse);
}
```

## **MsgPlaceBet**

Within this message, the user specifies the bet information they wish to place.

```proto
// MsgPlaceBet defines a message to place a bet with the given data
message MsgPlaceBet {
  // creator is the bettor address
  string creator = 1;

  // bet is the info of bet to place
  PlaceBetFields bet = 2;
}

// PlaceBetFields contains necessary fields which come in Place Bet Tx request
message PlaceBetFields {
  // uid is the unique uuid assigned to bet
  string uid = 1 [(gogoproto.customname) = "UID" ,(gogoproto.jsontag) = "uid", json_name = "uid"];

  // amount is the wager amount
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable)   = false];

  // ticket is a signed string containing important info such as `oddsValue`
  string ticket = 3;

  // odds_type is the type of odds bettor choose such as decimal, fraction
  sgenetwork.sge.bet.OddsType odds_type = 4;
}

// MsgPlaceBetResponse is the returning value in the response of MsgPlaceBet request
message MsgPlaceBetResponse {
    string error = 1;
    PlaceBetFields bet = 2;
}
```

### **Sample Place bet ticket**

```json
{
 "selected_odds": {
   "uid": "9991c60f-2025-48ce-ae79-1dc110f16990",
   "market_uid": "5531c60f-2025-48ce-ae79-1dc110f16000",
   "value": "2.0",
   "max_loss_multiplier": "1.0"
 },
 "kyc_data": {
   "ignore": false,
   "approved": true,
   "id": "sge1w77wnncp6w6llqt0ysgahpxjscg8wspw43jvtd"
 },
 "odds_type":1,
 "exp": 1667863498866062000,
 "iat": 1667827498,
 "iss": "Oracle",
 "sub": "CreateMarket"
}
```

### **Placement Failure cases**

The transaction will fail if:

- Basic validation fails:
  - Invalid creator address
  - Empty or invalid bet UID
  - Empty amount
  - Non positive amount
  - Empty or invalid ticket (containing space)
- Provided bet UID is already set
- Empty or invalid odds UID in ticket
- Empty, negative or invalid odds value in ticket
- Invalid bet value according to the selected `OddsType`
- There is no any market with the given marketUID
- The market is not active for accepting bet (it's not active or status in not `PENDING`)
- The market has expired
- The market maximum betting capacity has been reached
- The market does not contain the selected odds
- Bet amount is less than minimum allowed amount
- The creator address is not valid
- There is an error in AddPayoutProfitToMarket in market module
- There is an error in ProcessBetPlacement in Strategic Reserve module

### **What Happens if bet placement fails**

- The input data will not be stored in the `bet module's KVStore` and a meaningfull error will be returned to the client.

---

## **MsgSettleBet**

Within this message, the user provides a bet UID they wish to settle its corresponding bet.

```proto
// Msg defines the Msg service.
service Msg {
  // SettleBet defines a method to settle the given bet
  rpc SettleBet(MsgSettleBet) returns (MsgSettleBetResponse);
}
```

```proto
// MsgSettleBet defines a message to settle the given bet
message MsgSettleBet {
  // creator is the bettor address
  string creator = 1;

  // bet_uid is the unique uuid of the bet to settle
  string bet_uid = 2 [(gogoproto.customname) = "BetUID" ,(gogoproto.jsontag) = "bet_uid", json_name = "bet_uid"];

  // bettor_address is sthe bec32 address of the bettor
  string bettor_address = 3;
}

// MsgSettleBetResponse is the returning value in the response of MsgSettleBet request
message MsgSettleBetResponse {
    string error = 1;
    string bet_uid = 2 [(gogoproto.customname) = "BetUID" ,(gogoproto.jsontag) = "bet_uid", json_name = "bet_uid"];
}

```

### **Bet Settlement Failure cases**

The transaction will fail if:

- Basic validation fail:
  - Invalid creator address
  - Invalid bettor address
  - Empty bet UID
- Bet UID in invalid
- There is no matching bet for the bettor address
- Bet is canceled
- Bet is already settled
- Corresponding market not found
- Result of corresponding market is not declared
- There is an error in SR module functions

### **Settlement failure treatment**

- If corresponding market is aborted or canceled, the bet will be updated in the module state as below:

    ```go
    bet.Result = types.Bet_RESULT_REFUNDED
    bet.Status = types.Bet_STATUS_SETTLED
    ```

- Resolve the bet result based on the market result, and update field `Result` to indicate won or lost, and field `Status` to indicate result is declared.
- Call `Strategic Reserve module` to unlock fund and payout user based on the bet's result, and update the bet's `Status` field to indicate it is settled.
- Store the updated bet in the module state.

---
