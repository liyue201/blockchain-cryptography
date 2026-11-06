# The Four Properties of Hashing: The "Load-Bearing Wall" of Blockchain

An ordinary programmer uses hashing mostly for archiving or dedup. But in this chapter we want to make clear: hashing can **hold up an entire chain** because it possesses four never-negotiable properties at once. Remove one, and the chain collapses.

## ① Fixed Length
Whether you hash a single letter or an entire book, the output is always the same **fixed length** (SHA-256 is always 32 bytes / 256 bits).
- Benefit: uniform-sized "fingerprints", easy to compare, concatenate, and store.

## ② Avalanche Effect
If the input changes by even **one character**, the output changes drastically and looks completely unrelated.
- Benefit: whether something changed is instantly obvious (a single-byte comparison).

## ③ One-wayness
Given the output, **you cannot recover the input**; there is no shortcut other than brute force.
- Benefit: you can publish a hash without exposing the original (this is why "content → address" later can hide a private key).

## ④ Collision Resistance
It is practically impossible to find **two different inputs** that yield **the same hash**.
- Benefit: a hash can act as the content's "unique ID card".

> Run the code on the right and watch each of these four properties being demonstrated one by one. Focus especially on **[Avalanche]**: deleting just a single letter `a` changes nearly every one of the `32` bytes.