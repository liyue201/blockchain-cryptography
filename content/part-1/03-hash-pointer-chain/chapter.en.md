# Hash as Pointer: How Blocks Get "Chained" Together

With hashing that behaves like a "snowball", blockchain gets the most crucial word in its name — the **chain**.

## The Idea: Using a Hash as a "Hook"

Each block stores not only its own data, but also the **hash fingerprint of the previous block**. This "previous hash" works like a hook that firmly attaches the current block to the one before it — hence the name **hash pointer**.

```
[genesis] hash=A ─┬─► next
[block1] data + prev=A ─ hash=B ─┬─►
[block2] data + prev=B ─ hash=C
```

Because hashing has both the **avalanche effect** and **one-wayness**:
- if the data of any earlier block is changed, its hash changes;
- the "prev" stored in every later block then no longer matches → **the whole chain visibly breaks right before your eyes**.

## Why Is "Remembering One Hash" Enough?

You don't need to protect each block separately — **keep only the last hash of the chain**, and you can deduce whether the whole chain has been tampered with. This massively reduces "verification cost" — just as in real life, you only need to firmly memorize the ledger's "last page".

> Run the code on the right and focus on "after tampering block1": block1's hash changes, while block2's referenced `prev` is still the old value — they no longer match → the chain is caught red-handed.
> By hand change `blocks[1].data` to something clearly belonging to an "attacker", then run it, and feel how "one change, immediately exposed".