# **Events**

The OVM module emits the following events:

## *MsgSubmitPubkeysChangeProposalRequest*

| **Type**                    | **Attribute Key**         |      **Attribute Value**      |
|-----------------------------|---------------------------|-------------------------------|
| dvm_proposal_pubkeys_change | proposal_id               | {proposal_id}                 |
| dvm_proposal_pubkeys_change | content                   | {content}                     |
| message                     | module                    | market                        |
| message                     | action                    | dvm_proposal_pubkeys_change   |
| message                     | sender                    | {creator}                     |

## *MsgVotePubkeysChangeRequest*

| **Type**                    | **Attribute Key**         |      **Attribute Value**      |
|-----------------------------|---------------------------|-------------------------------|
| dvm_vote_pubkeys_change     | proposal_id               | {proposal_id}                 |
| dvm_vote_pubkeys_change     | voter_pubkey              | {voter_pubkey}                |
| dvm_vote_pubkeys_change     | vote                      | {vote        }                |
| message                     | module                    | market                        |
| message                     | action                    | dvm_vote_pubkeys_change       |
| message                     | sender                    | {creator}                     |
