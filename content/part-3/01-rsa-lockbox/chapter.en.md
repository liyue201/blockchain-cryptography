# The "Lockbox Model" of Public-Key Encryption: One Padlocks, the Other Opens

Finally we reach the core of Part 3 — **public-key (asymmetric) cryptography**. It solves the old headache of symmetric encryption: **how to safely hand the key to the other party**.

## The Lockbox Model (The Intuition Behind RSA)

Imagine a high-end safe that works with two keys:

```
Bob's public key  ➜  anyone can 'padlock' and drop something in
Bob's private key ➜  only Bob can 'open the box' and take it out
```

- **The public key can be public** (post it all over the street).
- **The private key stays absolutely secret**.
- **The directions are irreversible**: the key that locks can't unlock, and the key that unlocks can't lock. That's what "asymmetric" means.

## The Encryption Flow

1. Alice wants to send Bob a "secret"; she encrypts it with Bob's **public key** → padlocking it.
2. Only Bob's **private key** can decrypt → opening the box.
3. Even if others get the ciphertext, without the private key they can't open it.

## RSA's Real-World Limits

- It can only encrypt **small messages** (a 2048-bit key here can hold only ~190 bytes of plaintext);
- it's noticeably slower than AES;
- so real systems use "hybrid encryption": RSA wraps a **random session key**, which then symmetric-encrypts the real large data.

> Run the code on the right: watch anyone encrypt with the public key and only the private key can open it. **Don't worry yet about the math details**; first firmly grasp the direction of the "lockbox model" — it's the very same lock used two ways with the "digital signature" coming up next.