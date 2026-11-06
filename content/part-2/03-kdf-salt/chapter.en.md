# KDF & Salt: Protecting Your "Key Ring"

Wallet passphrases are usually "weak passwords" a human can remember. But encryption uses a 256-bit random key... so **how do you turn a "password" into a "reliable key"?** The answer is the **Key Derivation Function (KDF)**.

## The Old Problem: What's Wrong with Using the Password Directly as the Key?

- A password is a **human phrase**, far from a 256-bit random number, so its guessing space is small;
- An attacker can precompute "rainbow tables" of common passwords and their hashes to **look up them in an instant**.

## KDF's Three Moves

1. **Salt**: every user gets a **unique random salt**. Same password + different salt = different result, and rainbow tables collapse instantly.
2. **Slow (iterations)**: the derivation is deliberately made "computationally expensive" — PBKDF2 loops tens of thousands of times. The user only pays a tiny bit of slowness, but a cracker burns that much computing per password guess.
3. **Irreversible**: you can't recover the password from the key, so verification/authentication stays safe.

```
passphrase + unique salt + 10000 iterations ──► HMAC-SHA256 ──► 32-byte key
```

## What the Code Shows

- `pbkdf2Like(pass, salt, 10000, 32)`: a reimplementation of the PBKDF2-HMAC-SHA256 core (the real library such as `x/crypto/pbkdf2` is exactly this).
- Same password + salt ⇒ **stably yields the same key**.
- Changing one character / switching the salt ⇒ the key **changes beyond recognition**.

> Run the code on the right and focus on: with the same password and same salt, twice you get exactly the same 32-byte key; but changing the last character of the password immediately yields a different derived key.
> In production, Argon2/scrypt are also popular; the idea is the same: salting + slow iterations + memory-hardness.