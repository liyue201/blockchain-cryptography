# Merkle Proof: How a Light Node Verifies with Just One Path

In the previous section we packed transactions into a Merkle tree and obtained a "root". But wait — a **light wallet** (e.g. a mobile wallet) doesn't have that much storage; it can't download the whole blockchain. So how does it prove "my transaction is indeed in a certain block"?

The answer is the **Merkle proof**: you don't download the whole tree — you only carry a few "sibling hashes" along the path **from your leaf up to the root**.

## Imagine a Binary-Tree "Ladder"

Take 4 transactions as an example; the root is made of `H(topL, topR)`. To prove `Tx_B` is in the tree, the light wallet only needs to carry two hashes:
- `L0` (Tx_B's left sibling), and
- `topR` (the right half-subtree root of the side where Tx_B lives).

Then it "walks up" like this:
```
start from Tx_B
   combine with left sibling L0   → get topL
   combine with right sibling topR → get root
compare: is it the root the whole network publishes?
```

As long as it can reconstruct that root, it's 100% proof that `Tx_B` is in the tree — **even if you have no idea what C, D are**.

## Why Is This Remarkable?

- Data goes from O(n) down to **O(log n)**: for a block with 10,000 transactions, the proof needs only a dozen hashes.
- Given just one path, anyone can verify independently; the root is network consensus and cannot be forged.
- That's how a light wallet can "synchronize only the headers" while still securely confirming its own receipts.

> Run the code on the right: the target leaf `Tx_B` only carries the two "siblings" `L0` and `topR`, then reconstructs step by step from bottom to top back to `root`, matching what the network publishes.
> Try swapping out some value in the proof path (e.g. a fake sibling) — can it still match the root? If not, the proof immediately fails.