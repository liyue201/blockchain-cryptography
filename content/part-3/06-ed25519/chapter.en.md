# Ed25519: More Modern, Faster, and More Side-Channel-Resistant

Besides ECDSA + secp256k1, modern systems increasingly adopt **Ed25519**. It's not merely "another curve" — it's an upgrade of the whole design philosophy.

## Remember It in One Sentence

> Ed25519 handles, all at once, the three things most likely to go wrong in the ECDSA era — random-number reuse, side-channel leaks, and signature size.

| | ECDSA + secp256k1 | Ed25519 |
|---|---|---|
| random k | relies on the system RNG (reuse leaks the private key) | **deterministically derived** from (private key + message), no RNG source |
| side channels | demanding on implementation | constant-time, naturally resistant to timing attacks |
| signature length | 64~72 bytes | fixed 64 bytes |
| verify speed | fast | faster, simpler implementation |

## The Key: Making the "Random Number" Disappear

Ed25519 (a form of EdDSA) turns the random component needed for signing **into a deterministic function**, so:

- it no longer depends on the RNG's quality;
- the "reuse k → private key leaks" backdoor no longer exists;
- the same message always signs to the same signature, easy to cache and audit.

## Why Bother Learning It

- It's **beginner-friendly and extremely hard to implement wrong**;
- although on-chain secp256k1/ECDSA still dominates, Ed25519 appears more and more in cold wallets and interop protocols;
- future new chains increasingly favor it — it's a name worth knowing early in "chain cryptography".

> Run the code on the right: keygen → sign → verify passes → tampered fails, all in one flow. Note especially that the signature is a fixed 64 bytes and there's no "extra random number" anywhere.