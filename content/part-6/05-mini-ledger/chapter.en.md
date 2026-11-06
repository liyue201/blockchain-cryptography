# Finale: A Trustworthy Mini-Ledger (Integrated Lab)

Every part of the book has been delivered; this final section **assembles them into a runnable, trustworthy ledger**. On the right is a complete, compilable, runnable Go program that kneads all this math into a 3-block mini-chain:

| Part | Where it's used | Source part |
|---|---|---|
| hash-pointer chain | blocks interlock via prev; any change ripples through the chain | Part 6 §01 |
| Merkle root | folds a block's batch of transaction hashes into one "digest" | Part 1 |
| ECDSA signature | each transaction authorized by the sender's private key | Part 3 |
| nonce | transaction sequence, anti-replay | Part 3 / 6 |
| PoW hash puzzle | block hash must begin with several zeros | Part 6 §04 |

## What to Watch When Running

1. **Build**: three accounts Alice, Bob, Carol each hold a private key; `pay()` signs three transfers with the sender's private key, packed into 3 blocks with PoW difficulty;
2. **Verify**: `verify()` checks four things in one pass — hash-chain consistency, PoW prefix met, block linkage, and **per-transaction signature verification**; all pass returns `true`;
3. **Attack**: sneakily change one of `Alice`'s transfer amounts `→999999`; because the signature was made over the **pre-change** content, verification fails instantly and the whole chain returns `false`.

> This is the live version of "tampering is caught by math": the modifier lacks Alice's private key, so they can't produce a legitimate authorization for a transaction that doesn't exist.

## Why It's "Trustworthy"

- no need to trust a third-party "bookkeeper", because every node can **independently recompute and verify** the whole chain;
- without the private key, no one can alter any authorization — and if they do, it won't verify;
- to rewrite history you'd have to redo the PoW of every following block (an unbearable cost).

Security doesn't come from "whom you believe", but from **verifiable math** — the sentence running through this whole book finally finds its complete landing point here.

## Goodbye

By now you've packed the book's skeleton into your head: **one-way hashing, symmetric encryption, public/private keys & elliptic curves, digital signatures, ZKP foundations, multisig/aggregate/ring signatures, post-quantum directions**, and how they interlock into a trustworthy ledger. When you later read a real whitepaper, wallet, or smart contract, you'll recognize these "parts" again and again.

Blockchain cryptography — you haven't just "heard of it"; you've "calculated it through". Congratulations on finishing the whole journey.