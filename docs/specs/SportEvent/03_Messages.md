# **Messages**

The Sportevent module exposes the following services:

- AddEvent
- ResolveEvent
- UpdateEvent

```proto
// Msg defines the Msg service.
service Msg {
  rpc AddEvent(MsgAddEvent) returns (SportResponse);
  rpc ResolveEvent(MsgResolveEvent) returns (SportResponse);
  rpc UpdateEvent(MsgUpdateEvent) returns (SportResponse);
}
```

---

## **MsgAddEvent**

This message is used to add one or more new sportevent to the chain

```proto
message MsgAddEvent {
  string creator = 1;
  string ticket = 2;
}
```

---

## **MsgResolveEvent**

This message is used to resolve one or more already existent events on the chain

```proto
message MsgResolveEvent {
  string creator = 1;
  string ticket = 2;
}
```

---

## **MsgUpdateEvent**

This message is used to update one or more already existent events on the chain

```proto
message MsgResolveEvent {
  string creator = 1;
  string ticket = 2;
}
```

---

## **SportResponse**

This is the common response to all the messages

```proto
// common response for all the transactions call (batch transactions)
message SportResponse {
  string error = 1[(gogoproto.nullable) = true];
  SportEvent data = 2[(gogoproto.nullable) = true];
}

```
