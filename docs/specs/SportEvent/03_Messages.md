# **Messages**

The Sportevent module exposes the following services:

- AddEvent
- ResolveEvent
- UpdateEvent

```proto
// Msg defines the Msg service.
service Msg {
    rpc AddSportEvent(MsgAddSportEvent) returns (SportEventResponse);
    rpc ResolveSportEvent(MsgResolveSportEvent) returns (SportEventResponse);
    rpc UpdateSportEvent(MsgUpdateSportEvent) returns (SportEventResponse);
}
```

---

## **MsgAddSportEvent**

This message is used to add one or more new sportevent to the chain

```proto
message MsgAddSportEvent {
  string creator = 1;
  string ticket = 2;
}
```

---

## **MsgResolveSportEvent**

This message is used to resolve one or more already existent events on the chain

```proto
message MsgResolveSportEvent {
  string creator = 1;
  string ticket = 2;
}
```

---

## **MsgUpdateSportEvent**

This message is used to update one or more already existent events on the chain

```proto
message MsgUpdateSportEvent {
  string creator = 1;
  string ticket = 2;
}
```

---

## **SportEventResponse**

This is the common response to all the messages

```proto
// SportEvent response for all the transactions call
message SportEventResponse {
  string error = 1[(gogoproto.nullable) = true];
  SportEvent data = 2[(gogoproto.nullable) = true];
}
```
