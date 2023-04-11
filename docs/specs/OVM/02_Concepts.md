# **Concepts**

To verify the origin of the data, a simple mechanism of signing and verifiction with a private key and the correponding public key respectively is performed.

The data that needs to be pushed into the chain is first signed by the private key of a trusted source. This process occurs off chain. Essentially, we create an encrypted `ticket` with the trusted private key. The curve used for signing and verification of the ticket is the [Edwards-curve Digital Signature Algorithm](https://en.wikipedia.org/wiki/EdDSA). This algorithm has been chosen for the following benefits:

- Fast single-signature verification
- Fast batch verification
- Very fast signing
- Fast key generation
- High security level
- Foolproof session keys
- Collision resilience
- No secret array indices
- No secret branch conditions
- Small signatures
- Small keys

Details of the algorithm can be found [here](https://ed25519.cr.yp.to/)

---

After generating the encrypted signed `ticket`, this signature data is included in the transactions. This includes transactions for adding/editing betting markets on the chain, as well as verifying odds when the user places on these markets. All tickets come with an expiry timestamp which invalidates the ticket after a certain duration. This facility prevents the use or abuse of old and expired tickets.

The `OVM Module` essentially stores a list of trusted public keys. These public keys are just the counterpart to the private keys that were used to sign and encrypt the tickets off-chain. When a transaction is made to the chain that necessitates verification of the origin of the data, the corresponding module invokes the OVM module for verification and decryption purpose. The OVM Module works as an interface, which can decode any signed data passed to it when supplied with the encryption algorithm and the decrypted type. This design completely nullifies the need to change the structure of the OVM if the ticket structure changes. This essentially enables the OVM to be a global verification module.

When the OVM is invoked, it first attempts to verify the signature of the data against the list of registered public keys. If the signature is verified successfully by the leader public key (the first element in the public keys slice in the key vault store), the OVM decrypts the data into the provided structure and returns it to the invoking module. In case the signature verification fails, or the ticket seems to have expired, or the decrypted structure does not match with the expected structure, the verification is considered to be a failure and a corresponding error is returned to the invoking module, which consequently results in failure of the transaction.

The OVM Leader is the public key that is being used for the verification of the tickets, The first element of the Key Vault public keys is the leader. If the private key of the leader's public key gets corrupt/hacked/leaked, The holders of the rest of the public keys can create a pub keys-change proposal to replace the leaked public key with a new one and choose the new leader key. Each proposal needs at least 66.67% "yes" or 2 "no" to make blockchain code to decide about the proposal's approval or rejection, the modification happens in the end blocker if the condition of votes is satisfied.
