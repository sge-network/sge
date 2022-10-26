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

After generating the encrypted signed `ticket`, this signature data is included in the transactions. This includes transaction for adding/editing betting events on the chain, as well as verifying odds when user places on these events. All tickets come with an expiry timestamp which invalidates the ticket after a certain duration. This facility prevents use or abuse of old and expired tickets.

The `DVM Module` essentially stores a list of trusted public keys. These public keys are just the counterpart to the privatekeys that were used to sign and encrypt the tickets off-chain. When a transaction is made to the chain that necessitates verification of the origin of the data, the corresponding module invokes the DVM module for the verification and decryption purpose. The DVM Module works as an interface, which can decode any signed data passed to it when supplied with the encryption algorithm and the decrypted type. This design completely nullifes the need to change the structue of the DVM if the ticket structure changes. This essentially enables the DVM to be a global verification module.

When the DVM is invoked, it first attempts to verify the signatue of the data against the list of registered public keys. If the signature is verified successfully by at least one public key, the DVM decrypts the data into provided structre and return it to the invoking module. In case the signature verification fails, or the ticket seems to have expired, or the decrypted structure does not match with the expected structure, the verification is considered to be a failure and a corresponding error is returned to the invoking module, which consequently results in failure of the transaction.
