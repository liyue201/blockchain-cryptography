# The Nature of an Address: Deriving Public Key → Hash

Many think "an address is where the account balance lives". More precisely: **an address is a fingerprint (hash) of the public key, and asset ownership is decided by the private key's signing power behind it.**

## The Derivation (recap of Part 4 §02)

```
private key --(scalar mult)--> public key --(hash)--> pubkey hash (fingerprint) --(version+checksum+Base58)--> address
```

- **private key → public key**: elliptic-curve scalar multiplication (one-way);
- **public key → fingerprint**: hash condensed to a fixed length (one-way);
- **in reality one account can rotate through countless addresses**: all are different "faces" of the same public key (version, truncation, derivation path, etc.).

> Run the code on the right: one key first yields a 65-byte public key, then through hashing a 32-byte fingerprint, and can even derive "another address". Behind them all is the same public key, just presented differently externally.

## Why an "Address" Can Masquerade as an "Account"

- to the user: an alphanumeric string is far friendlier to copy than a 65B hash;
- to the chain: an address is just an index; balance/ownership is established by **private key → signature verification** (see next section);
- so "address = account" is a practical illusion — the truth is "**address = public-key fingerprint; ownership = private key's signing power**".

## Echoing the Whole Book

An address's birth simultaneously uses Part 1's hashing, Part 3's elliptic curves, and Part 4's encoding & checksum — this section is where they're twisted into the "nameplate to the outside world". Next section we put it to use: **authorize a transaction (sign) with the private key, and use nonce for anti-replay**, so "the money in this address" truly flows as you wish.