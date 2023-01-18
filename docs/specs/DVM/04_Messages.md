# **Messages**

In this section, we describe the processing of the DVM messages.

## **MsgMutation**

```proto
// Msg defines the Msg service.
service Msg {
      rpc Mutation(MsgMutation) returns (MsgMutationResponse);
}

message MsgMutation {
  string creator = 1;
  string txs = 2;
}

message MsgMutationResponse {
  bool success = 1;
}
```

This message is used to add or delete the trusted public keys to the DVM. In the absence of any keys, there is no verification for adding a new key to the DVM. However, once at least one key exsts, for adding or deleting keys, the signature needs to be verified by at least one already existing public key.

> **NOTE:** In the absence of public keys, signatures cannot be verified for any transaction.
