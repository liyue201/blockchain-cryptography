# Threshold & Multisig: How Multi-Sig Wallets Work

One person holding all the private keys is dangerous: if stolen, everything's gone; if lost, you can never get back in. **Multisig (multiple signature)** splits the power into pieces — **a certain number of people must be gathered** to move funds, i.e. "one short of the team means it can't be done".

## Meaning and Rules

- `n`-of-`m`: with `m` total members, at least `n` signatures from the "roster" are required to release;
- **2-of-3**: a company/family jointly manages — two can move funds, one being unavailable doesn't block, one turning traitor still can't drain it;
- this is mainstream in custody transfers, smart-contract wallets, and DAO treasuries.

> Run the code on the right: build 3 directors and try — **only 1 signature: rejected**; **directors A + C signing together: 2/3 passes, released**. Multisig's "threshold authorization" is crystal clear.

## Two Implementation Layers

1. **On-chain multisig (straightforward)**: the contract registers several public keys and a threshold `n`; each payment executes once ≥ n signatures are collected. Simple and reliable.
2. **Threshold / aggregate signatures (better)**: the participating `n` people each sign, and the result is assembled into **one** signature representing the whole group (e.g. MuSig, BLS — the next section's BLS aggregation); the chain sees one signature, one address, and can't see how many or which people are behind it. Saves space and protects participant privacy.

## Which Risks It Covers

- **Against single-point theft**: stealing one director's private key can't gather enough signatures — useless;
- **Against single-point failure**: losing one director's key, the rest still operate (as long as ≥ n);
- **Governance**: multisig is the cryptographic incarnation of DAO "voting = spending" — you vote with signatures, not with your mouth.

> Summary: multisig turns "the private key is everything" into "many minds pulling toward one goal"; threshold/aggregate signatures then make it maximally lean. This section and next section's BLS aggregation are the union of "more secure" and "smaller".