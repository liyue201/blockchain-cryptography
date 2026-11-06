# Elliptic Curve Cryptography (ECC) Intro: A Curve Full of "One-Way Secrets"

In sections 1 and 2 we used RSA to build a sense of direction. But what **blockchain actually uses is elliptic curves (ECC)**. In this section we set details aside and grasp only the core intuition.

## A Curve and a Group Operation

An elliptic curve is a set of points satisfying `y² = x³ + a·x + b`. The key: **two points can be "added", and a point can "add itself"** — forming a closed "group operation".

And the truly important operation is **scalar multiplication**:

```
k·G  = G + G + ... + G   (k point additions)
```

## One-wayness: ECC's Security Foundation

Like RSA, ECC rests on an asymmetric capability:

- **Forward (easy)**: given the algebra k and the base point G, computing `k·G` takes a single scalar multiplication — done in a blink by the CPU;
- **Reverse (extremely hard)**: given `k·G`, recovering k requires solving the "elliptic curve discrete logarithm problem" — there's no known polynomial-time algorithm, and the 2^256 space can't be exhausted.

**That's where the private/public key pair comes from:**
```
private key k ──scalar-mult──► public key Q = k·G
reverse Q → k   「impossible」
```

## Why It Beats RSA

At the same security strength, ECC's keys are much shorter (256 bits ≈ RSA 3072-bit or more), signatures are smaller and verification faster — **every transaction's signature on the chain must be re-verified by every node in the world**, and by that accounting ECC wins outright.

> Run the code on the right and watch a random private key k turn into the public key k·G, noting that `2·G` is also a completely different point. Truly understand: forward is instant, reverse never achievable.
> I'm not trying to make you draw the curve beautifully — remember that summary line — **"from k to k·G is fast; the reverse has no road"**.