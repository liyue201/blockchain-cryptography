# A Bird's-Eye View of Blockchain from a Cryptographic Lens

The previous five parts covered the **parts**; this part **assembles them into a chain**. Let's start from the most fundamental skeleton: **blockchain = a chain threaded together with hash pointers**.

## Hash-Pointer Chain: How Blocks Lock Into Each Other

Each block contains at least `data` and a `prevHash` — **the hash of the previous block** (a hash pointer). Thus:

- block `i`'s hash embeds the hash of block `i-1`;
- block `i+1`'s hash in turn embeds block `i`'s hash…

A **chain of interlocking links**: if any early block's content changes, its hash changes, dragging every subsequent block's hash with it.

> Run the code on the right: generate 5 blocks and see that each block's hash indeed contains the previous block's prevHash (`#1`'s prevHash = `#0`'s hash …). `verifyChain` reports `chain integrity: true`; when someone changes the 1st block's `A→B5` to `A→B999`, `verifyChain` flips to **false** immediately — the change dominoes out to the end of the chain.

## Why It's "Tamper-Proof"

Not "others can't modify the file", but "modifications are instantly detected":

- change one spot → that block's hash changes → mismatches the next block's prevHash → the chain breaks;
- to forge, you'd have to **recompute every subsequent block** from the change point (huge work, especially combined with Part 4's PoW);
- and the chain isn't held by any single party — the whole network has a copy, so fooling everyone at once is nearly impossible.

## Summing Up the Book in One Sentence

Chapter one taught "hashes are one-way"; this chapter turns it into "an interlocking ledger". **Blockchain converts the human-societal problem of "trust" into a "verifiable mathematical problem"** — that's cryptography's most fundamental value on-chain.

> Next section drops **public-key hash → address** onto the ledger, and then assembles **signature authorization, PoW hash-finding, and anti-replay** one by one into a runnable mini-ledger.