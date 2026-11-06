# Proof of Work (PoW): A Cryptographic Hash Puzzle

Hashes can do more than "compress into a fingerprint" — they can serve as "proof of labor". **PoW (Proof of Work)** at its essence: find a hash that **starts with `difficulty` zeros**; hardness rises exponentially with the prefix length.

## Why "finding the nonce" Is So Hard

- **hard to find**: to get the first `d` digits of the hash all zero, you must keep changing the `nonce` and trying one by one; each candidate hash is like rolling dice, with success probability `1/16^d`;
- **easy to verify**: given `(nonce, hash)`, anyone recomputes the hash once to confirm instantly — hard to find, nearly free to verify.

> Run the code on the right: difficulty grows from 3 to 5, and attempts jump from tens of thousands to over a million. **Difficulty +1 ⇒ expected work ×~16.** Watch the same `data`: the larger the nonce, the more the hash's leading segment looks like pure luck.

## What It Adds to the Chain

1. **makes "writing a block" expensive**: to forge a long chain you must recompute every later block's PoW, and the cost quickly becomes unbearable;
2. **forms a consensus anchor**: the chain able to accumulate the most PoW (the longest chain) is deemed most trustworthy — nudging the whole network to follow it;
3. **remember the key point**: here PoW is just a **hash puzzle** (a cryptographic concept); the real "consensus/network sync" is a separate systems-engineering topic — this book only covers the cryptographic part it depends on.

> Echoing Part 1 §06 (the cryptographic essence of mining): back then we saw its power with a hand-written nonce; now that it's back in the larger context of "making the chain tamper-proof", you see why the hash puzzle is the chain's "manual-labour endorsement".