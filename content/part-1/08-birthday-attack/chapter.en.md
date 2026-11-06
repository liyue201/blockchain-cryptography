# The Security Threshold of Hashing: Birthday Attack & Collision Resistance

We keep stressing "hashes must be long enough to be safe". This section explains **why**, and the answer hides in a counter-intuitive mathematical phenomenon — the **birthday paradox**.

## The Birthday Paradox: Why Collisions Are Easier to Find Than You'd Think

In a class of only ~23 people, there is a better-than-50% chance that **two people share a birthday**. The expected number for "one collision" isn't 365, but about **√365 ≈ 19**. Translated into hash terms:

> To find a hash collision where `the first n bits match`, you only need to try about **2^(n/2)** attempts, not 2^n.

That's exactly the "birthday attack": the attacker doesn't hunt for **some specific bits** matching — they keep generating many messages, waiting for any **two arbitrary ones** to collide.

## What Does This Mean for Choosing a Size?

| Desired security level | Required hash size |
|----------------------|--------------------|
| Colliding via the birthday bound 2^(n/2) | n must be large enough |
| SHA-1 (160 bits) | birthday bound 2^80; long believed safe in theory, but collisions were ultimately found |
| SHA-256 (256 bits) | birthday bound 2^128, unreachable by every current machine → **safe** |
| Achievable scale | modern chips compute roughly 2^70 ops/sec-class; 2^128 is astronomical |

That's why the industry treats **256 bits** as the "minimum safe line": it pushes collision difficulty beyond anything engineering can reach.

> Run the code on the right and see with your own eyes: finding a collision where "the first 16 bits match" needs only ~256 attempts (≈ 2^(16/2)), whereas intuitively you might expect 65,536. That explains why hashes need to be as long as 256 bits.