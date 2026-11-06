# Zero-Knowledge Proofs (ZKP): Proving Without Revealing

**Zero-knowledge proofs** are cryptography's most fascinating achievement: **I can prove "I know something" without revealing a single word about that thing itself.** They are the foundation of privacy coins, zk-rollups, and cross-chain bridges.

This book won't write out full ZK math (that needs pairings and polynomial commitments); instead we use a **hash commitment** to give you ZKP's first-principles intuition — because all modern ZK is built on "commitments".

## What a Commitment Is

A commitment = **lock it up first to be accountable, open it later to show truth**, in two phases:

- **Commit phase**: the Prover picks a secret `secret`, generates a random `nonce`, and publishes `commit = H(secret ‖ nonce)` (only the hash, not the content);
- **Open phase**: the Prover reveals `(secret, nonce)`, and the Verifier recomputes the hash to compare; if it matches, it's accepted.

> Run the code on the right: the Prover publishes a commitment `c`; the Verifier can't guess the secret from `c` (**hiding**); once opened, comparison matches immediately; and if the Prover tries to "change their story", the hash won't match so they **can't deny it (binding)**.

## The Two Security Properties of Commitments (Threading the Book)

1. **Hiding**: seeing `c = H(secret ‖ nonce)` reveals nothing about `secret`. The reason is Part 1's **one-wayness** — especially since that random `nonce` lets the same secret yield completely different commitments.
2. **Binding**: once the commitment is given, the Prover cannot open it to the same `c` with a different `(secret', nonce')` (that would be a hash collision). It locks in "no going back".

## From Commitment to Zero-Knowledge

The real not-revealing part: **the Prover shows only a verifiable "relationship holds", not the secret itself.** A commitment turns the "secret" into "I've locked up a secret"; ZK then makes "open/verify" need no full reveal — merely proving "an openable answer exists".

So privacy coins (Monero/Zcash) use ZK to achieve "public ledger, private amounts & identity"; zk-rollups compress tens of thousands of transactions into one one-click-verifiable proof. **Public and private, both attainable thanks to cryptography.**

> ⚠️ Note: this example is "commit/open", an introduction to ZKP; full ZK also needs "proving validity without opening" (e.g. Groth16/PLONK), which is frontier beyond Part 5 — but the foundation lies in these few dozen lines you just ran.