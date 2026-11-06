# ECDSA: The Most Mainstream Signature Algorithm on Chains

Bitcoin and Ethereum transactions are signed with **ECDSA (Elliptic Curve Digital Signature Algorithm)**. It puts "sign with private key, verify with public key" onto the secp256k1 curve.

## Reading the Signature Formula in One Go

Let the message hash be `z`, the private key `d`, a random `k`, and the base point `G`:

```
r = (k·G).x  mod n
s = k⁻¹·(z + r·d)  mod n
signature = (r, s)
```

Verification only needs the public key `Q = d·G`:

```
u₁ = z·s⁻¹,  u₂ = r·s⁻¹
does (u₁·G + u₂·Q).x  mod n  equal r ?
```

**Why does it hold?** `u₁·G + u₂·Q = (z·s⁻¹ + r·s⁻¹·d)·G = s⁻¹(z + r·d)·G = k·G`, landing exactly back at the x-coordinate r. That algebraic identity is the entire secret behind "signatures being verifiable".

## That Ubiquitous "Random Number k"

It's powerful, but also risky — **k must never leak or be reused**! Because once k leaks, the private key can be immediately solved as `d = (k·s - z)·r⁻¹`; over the years there have been multiple "random-number reuse" incidents where huge private keys were computed this way.

## The Fix: RFC6979 Deterministic Signing

Derive the **unique k from the private key and the message hash** (rather than relying on the system RNG), and signing becomes "deterministic":

- same signer + same message = **always the same signature**, eliminating k reuse;
- more stable, reproducible, and easier to audit.

> Run the code on the right to witness three things:
> ① the hand-written ECDSA sign/verify math fully works (original message passes, tampered fails);
> ② with an RFC6979-style deterministic k, **signing the same message twice gives byte-identical signatures**;
> ③ finally the Go stdlib `ecdsa.Verify` cross-checks our hand-written signature — valid.
> To feel the "price of k leaking", swap `deterministicK` for a fixed random number and sense the security risk level — it's right here that you grasp "the random number IS the private key".