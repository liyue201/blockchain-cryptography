# Getting Random Numbers Right: The Most Expensive Mistake On-Chain

Part 3 mentioned the random number `k` when covering ECDSA (`(r,s) = (k·G, k⁻¹(z + r·d))`). You might think k is just "a little helper". **This section tells you: use k wrong and your private key goes to zero on the spot.** Many historic coin losses — especially hacked Bitcoin/Ethereum arbitrage bots — were k reuse.

## The Essence of the Attack: k Leak = Private Key Leak

If the same private key signs **two different messages**, both times using the same `k`:

- since `r = (k·G).x`, **same k → r must be identical** (a flaw visible to the naked eye);
- the two signatures give `s1 = k⁻¹(z1 + r·d)` and `s2 = k⁻¹(z2 + r·d)`;
- **subtracting the two equations** cancels the private-key term; first solve for `k`, then substitute back into either to solve for `d`:

```
k = (z1 - z2) / (s1 - s2)        (mod n)
d = (s1·k - z1) / r              (mod n)
```

> Run the code on the right: with two signatures sharing the same `r`, the "thief" precisely recovers the victim's true private key `d` (byte-identical). This isn't armchair math — it's a real attack you can run in a few dozen lines.

## Why It's Still So Easily Triggered

- someone lazily sets k to a **constant** → guaranteed reuse;
- a broken RNG (not enough entropy, reused seed — see Part 4 §01/02) produces the same k;
- some implementations have predictable `k` across platforms, reconstructed by academics with algorithms.

## The Right Approach (modern wallet's standard answer)

1. **RFC6979**: **deterministically derive** k from `(private key, message hash)` — the same message always gets the same k, but each message naturally differs, eliminating the "random" step and hence reuse;
2. **Ed25519** (Part 3 §06) directly **derives k by hashing private key + message**, removing the "random k" variable by design;
3. never hand-write "random k" — leave it to audited libraries.

> This section and 41/42 are two sides of one theme: **the private key itself must be truly random (entropy); the ephemeral k at signing time must be even more rigorous.** Using randomness wrong is the most expensive mistake on-chain — and also among the easiest to prevent.