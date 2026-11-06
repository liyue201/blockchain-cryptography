# Comparing Mainstream Hash Algorithms: Why Chains Choose SHA-256

"Hash" is the umbrella term for a whole class of methods; the specific "recipes" come in many varieties, like a row of different "fingerprint machines". The common ones: **MD5, SHA-1, SHA-256, SHA-512**.

## A Row of "Fingerprint Machines": the Longer the Output, the Harder to Collide

| Algorithm | Output length | Status | Notes |
|-----------|---------------|--------|-------|
| **MD5** | 128 bits | ❌ broken | colliding since 2004, only usable for non-security checks |
| **SHA-1** | 160 bits | ❌ broken | SHAttered collision demo in 2017; must be retired |
| **SHA-256** | 256 bits | ✅ secure | the de-facto standard of the blockchain world (widely used by Bitcoin/Ethereum) |
| **SHA-512** | 512 bits | ✅ secure | longer output for higher security needs, but a bit slower |

## A Key Idea: Output Size = Security Margin

Why is "longer = safer"? Because it directly determines **collision difficulty**:
- 128-bit output: colliding theoretically needs about 2^64 operations — **reachable** by modern computers → dangerous.
- 256-bit output: needs about 2^128 operations — an **astronomical number** beyond every known computer → safe.

## So the Answer Is…

Blockchains choose **SHA-256** because:
1. 256-bit output gives enough collision-resistance margin;
2. it has all four properties — **fixed length, speed, avalanche, one-wayness**;
3. hardware (ASIC miners) ecosystem is mature and the implementation has been battle-tested at scale.

> Run the code on the right, feeding the same sentence to each of the four "fingerprint machines", and watch the output lengths (128/160/256/512 bits) differ.
> Try changing any single character and see whether SHA-256 "avalanches" again.