# **Events**

The OVM module emits the following events:

## *MsgSubmitPubkeysChangeProposalRequest*

| **Type**                    | **Attribute Key**         |      **Attribute Value**      |
|-----------------------------|---------------------------|-------------------------------|
| ovm_proposal_pubkeys_change | proposal_id               | {proposal_id}                 |
| ovm_proposal_pubkeys_change | content                   | {content}                     |
| message                     | module                    | market                        |
| message                     | action                    | ovm_proposal_pubkeys_change   |
| message                     | sender                    | {creator}                     |

## *MsgVotePubkeysChangeRequest*

| **Type**                    | **Attribute Key**         |      **Attribute Value**      |
|-----------------------------|---------------------------|-------------------------------|
| ovm_vote_pubkeys_change     | proposal_id               | {proposal_id}                 |
| ovm_vote_pubkeys_change     | voter_pubkey              | {voter_pubkey}                |
| ovm_vote_pubkeys_change     | vote                      | {vote        }                |
| message                     | module                    | market                        |
| message                     | action                    | ovm_vote_pubkeys_change       |
| message                     | sender                    | {creator}                     |
