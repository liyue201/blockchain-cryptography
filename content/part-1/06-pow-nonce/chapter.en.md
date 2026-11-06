# The Cryptographic Essence of Mining: A Hash Puzzle with Tunable Difficulty

You often hear "mining", "hashrate", "difficulty". Setting the hype aside, from a cryptography standpoint it is actually a very plain **hash puzzle**.

## What the Puzzle Looks Like

> Find a number `nonce` such that `SHA-256(block header + nonce)` starts with a certain number of `0`s.

- **Easy to verify**: pick any nonce, hash it once, and you immediately know yes or no.
- **Hard to find**: because hashing is "unordered", there is no trick — you can only **try one by one**. That's the "work".

## Tunable Difficulty: How Hard You Want It, How Long It Takes

The more leading `0`s you require, the higher the difficulty:
- require the first 4 bits (1 hex) to be 0: average 2^4 = 16 tries.
- require the first 8 bits (2 hex) to be 0: average 2^8 = 256 tries.
- require the first 16 bits: average 2^16 = 65,536 tries.
- require the first 128 bits? That would need 2^128 tries — not even a destroyed universe could finish.

Each extra `0` multiplies the work by about **×16**. A whole network can use one difficulty value to "tune" the block time near a target (e.g. Bitcoin's ~10 minutes).

## Why Does This Hold in a "Chain"?

Because hashing is **one-way + shortcut-free**: apart from brute-forcing the nonce, there is no way to know in advance which nonce will win. The "hashrate" miners invest is essentially the computation burned on a shortcut-free problem — which makes "arbitrarily rewriting the chain" extremely expensive.

> ⚠️ Note: here we only cover the cryptographic fact of the "hash puzzle". How PoW organizes consensus, the longest-chain rule, etc., belongs to blockchain/consensus-layer knowledge, which is beyond this book's scope.

> Run the code on the right to experience the nonce-finding speed under different difficulties, and intuitively feel "each extra 0 makes the work ×16".