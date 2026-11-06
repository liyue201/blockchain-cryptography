# Commitment Schemes and the Random Oracle

Last section we made a commitment with a hash. This section digs into the **two security properties** of commitments and the backstage character behind them — the **random oracle**.

## The Random Oracle

In the real world there's no "true random function"; we use a hypothetical model, the **random oracle**, to describe an ideal hash: **any input yields a fixed but seemingly random, unpredictable output.**

A cryptographic hash (SHA-256) is widely regarded as "playing the role well enough":

- **deterministic**: same input always gives the same output (keeping verification reproducible);
- **unpredictable / avalanche**: a one-bit difference in input, yet the output is like rolling dice (making reversal impossible).

> Run the code on the right: changing the nonce for the same secret changes the commitment completely; opening with a different secret fails to match; then list a few inputs and see how "patternless" the hash outputs are. Those three snippets are direct demonstrations of commitment properties and RO behavior.

## Commitment = Two Gates of a Harbour

| Property | Meaning | Which RO trait it relies on |
|---|---|---|
| **Hiding** | seeing only `c` reveals nothing about the secret | one-way irreversibility |
| **Binding** | content can't be changed after committing | deterministic output, collision-free |

The two gates together form the game rule of "honestly abiding by a commitment": **before I open, you can't guess; once I open, you can't deny.**

## Why Chains Depend on It So Much

- **Commit-then-reveal**: parties first each "lock" their answer (commitment), and only open after all are submitted — preventing someone from watching others' answers and changing theirs (auctions, lotteries, random-number generation all do this);
- **Fiat–Shamir transform**: turning a ZK proof non-interactive works precisely by replacing "the verifier's random challenge" with "a hash over the message" — all running on the random oracle.

Understanding commitments and the RO hands you the key from "one-to-one proofs" to "proofs on a public ledger (zk-rollup / privacy chains)".