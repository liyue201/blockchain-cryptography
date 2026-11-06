# From Content to Address: A "Derivation Chain" of Hashes

On the chain, to send money to someone you need an **address**. An address isn't a randomly chosen name — it is "derived" from your **public key through a series of hashes**. That's a watertight derivation chain:

```
private key ──► public key(X,Y) ──► public key hash ──► address
                 elliptic curve            hash          encoding
```

## Step by Step (matching the code on the right)

1. **Private key**: a 256-bit random number. Whoever has it owns the money.
2. **Public key**: a point on the elliptic curve computed from the private key. It can't be reversed into the private key (one-way).
3. **Public key hash `hash160`**: hash the public key to get a 160-bit (20-byte) fingerprint.
   - Real chains use `RIPEMD160(SHA256(pubkey))`; this example uses double `SHA-256` as a teaching approximation.
4. **Checksum**: the first 4 bytes of a double SHA-256, appended at the end, to **detect typos in the address**.
5. **Base58Check encoding**: translate that binary block into a human-friendly "address string" that omits confusing characters like `0 O I l`.

## Why the Extra "Detour Through a Hash"?

- **Hides the public key**: others only see your "hashed address" and can't reverse it into public key → private key.
- **Fault tolerance**: even a single wrong character usually breaks the checksum, so the system can instantly flag "invalid address".
- **Privacy/security**: one more one-way hash is one more wall.

> Run the code on the right and watch the full chain: private key → compressed public key → 20-byte hash → payload with checksum.
> Note the `checksum`: change just one hex character and the checksum no longer matches — that matches the fault-tolerance design of the real world.