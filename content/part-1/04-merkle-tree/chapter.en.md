# The Merkle Tree: Rolling Up a Pile of Hashes into One Root

A single chain can hold **hundreds of thousands of transactions**. But if "verifying the whole chain" required comparing every transaction one by one, the cost would be unbearable. The Merkle tree is the big weapon that answers this: **compress a pile of transactions into one "root", and if any single transaction changes, the root changes too.**

## How Is It "Stacked Up"?

Hash each transaction first, getting a row of **leaves**; then **pair them up and hash again**, layer by layer, until only one remains — called the **Merkle root**.

```
         root = H( mid0 + mid1 )
              /                \
      mid0=H(L0,L1)         mid1=H(L2,L3)
      /        \            /        \
   L0=H(Tx_A) L1=H(Tx_B)  L2=H(Tx_C) L3=H(Tx_D)
```

## This "Tree" Brings Two Astonishing Benefits

1. **One root covers the whole ledger**: the block header only needs to record a single `merkle_root` to represent every transaction in that block.
2. **Change once, exposed immediately**: if any transaction changes, it propagates up along its own path to the root, and once the root changes the whole network notices.

> Run the code on the right and watch how 4 transactions gradually merge into a single Merkle root; then change any transaction to `Tx_X` and observe the **root changing immediately**.

The beauty of the Merkle tree doesn't end here — in the next section we'll see how a light wallet can prove "my transaction is indeed in the block" by **downloading only one path**, without downloading the whole block.