# Why Chains Favor secp256k1: Curve Parameters & Selection

Open the Bitcoin/Ethereum source code and you'll almost always see this tongue-twister of a name — **secp256k1**. What's so special about it?

## A "Minimalist" Curve

Most elliptic curves tweak parameters to be more complex for performance/security (e.g. P-256). secp256k1 goes the opposite way:

```
y² = x³ + 7     ← just x³ plus 7, coefficients a=0, b=7
```

- **Coefficients are only 0 and 7**: point addition/doubling is extremely cheap, the code is small, and there's less room for hidden backdoors;
- **A solid base**: there's also a carefully chosen base point G with a large prime-order n, 256 bits long;
- **Plenty of security margin**: private-key space 2^256, birthday bound 2^128 — nothing currently in service can touch it.

## Selection Differences: A One-Line Takeaway

| Curve | Trait | Common use |
|-------|-------|-----------|
| secp256k1 | simplest coefficients, community heavily audited | Bitcoin / Ethereum mainnet |
| P-256 (secp256r1) | chosen by standards bodies, many reference implementations | TLS site certs and other general-purpose use |
| Ed25519 | more modern, constant-time against side channels | private-key signing in new systems |

All three deliver security at the 256-bit level; the differences chiefly lie in "how the parameters were chosen and how battle-tested the implementation is". The chain world settled on secp256k1 by decades of convention, so network-wide verification logic stays unified.

> Run the code on the right: this example needs no third-party library — it directly uses Go's `elliptic.CurveParams` to **manually construct secp256k1's public domain parameters** (the prime p, a=0, b=7, base point G, order n), then runs private key → public key once more.
> In other words, secp256k1's "whole recipe" is just those strings of public numbers, reproducible by anyone — that's precisely where its auditability lies.