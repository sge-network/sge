# **Messages**

In this section, we describe the processing of the OVM messages.

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

This message can be sent by any of the current registered public key owners, If any of the current private keys get compromised, we can use this message to replace the corrupt public key and set the new leader.

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

### **PublicKeysChangeProposalPayload**

```proto
// PubkeysChangeProposalPayload indicates data of public keys changes proposal
// ticket.
message PubkeysChangeProposalPayload {
  // public_keys contain new pub keys to be added to public keys.
  repeated string public_keys = 1;
  // leader_index is the universal unique identifier of the public key.
  uint32 leader_index = 2;
}
```

## **VotePubbkeysChange**

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
  // voter_key_index is the public key index of the voter in the current list
  // of public keys in the vault.
  uint32 voter_key_index = 3;
}
```

```proto
// MsgVotePubkeysChangeResponse is the type of response vote for public keys
// modification.
message MsgVotePubkeysChangeResponse { bool success = 1; }
```

### **VotePubkeysChangePayload**

```proto
// ProposalVotePayload indicates vote data ticket.
message ProposalVotePayload {
  // proposal_id is the id of the proposal.
  uint64 proposal_id = 1;
  // vote is the chosen option for the vote.
  ProposalVote vote = 2;
}
```

> **NOTE:** In the absence of public keys, signatures cannot be verified for any transaction.
