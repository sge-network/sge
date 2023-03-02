# **Messages**

In this section, we describe the processing of the DVM messages.

## **SubmitPubkeysChangeProposal**

This endpoint accepts the message for creating a proposal for the pubkeys change.

```proto
service Msg {
  // PubkeysChangeProposal defines a method to submit a proposal for changing
  // the allowed public keys.
  rpc SubmitPubkeysChangeProposal(MsgSubmitPubkeysChangeProposalRequest)
      returns (MsgSubmitPubkeysChangeProposalResponse);
}
```

### **MsgSubmitPubkeysChangeProposalRequest**

```proto
// MsgPubkeysChangeProposalRequest is the type of request for modification of
// public keys.
message MsgSubmitPubkeysChangeProposalRequest {
  // creator is the account address of the creator.
  string creator = 1;
  // ticket is the jwt ticket data.
  string ticket = 2;
}
```

```proto
// MsgPubkeysChangeProposalResponse is the type of response for modification of
// public keys.
message MsgSubmitPubkeysChangeProposalResponse { bool success = 1; }
```

## **VaotePubbkeysChange**

This endpoint accepts the message to vote on a proposal for the pubkeys change.

```proto
service Msg {
  // VotePubkeysChange defines a method to vote for a proposal for changing the
  // allowed public keys.
  rpc VotePubkeysChange(MsgVotePubkeysChangeRequest)
      returns (MsgVotePubkeysChangeResponse);
}
```

### **MsgVotePubkeysChangeRequest**

```proto
// MsgVotePubkeysChangeRequest is the type of request to vote on the
// modification of public keys proposal.
message MsgVotePubkeysChangeRequest {
  // creator is the account address of the creator.
  string creator = 1;
  // ticket is the jwt ticket data.
  string ticket = 2;
  // public_key is the public key of the voter.
  string public_key = 3;
}
```

```proto
// MsgVotePubkeysChangeResponse is the type of response vote for public keys
// modification.
message MsgVotePubkeysChangeResponse { bool success = 1; }
```

> **NOTE:** In the absence of public keys, signatures cannot be verified for any transaction.
