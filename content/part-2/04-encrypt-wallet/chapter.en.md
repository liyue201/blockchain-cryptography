# Where Symmetric Encryption Shines Off-Chain: Encrypting Your Wallet

After hashing and signatures, you might ask: **where does symmetric encryption actually get used on-chain?** Straight to the point: it rarely appears "on-chain", but it's **everywhere off-chain** — the most typical case is **encrypting the wallet file**.

## What Does a Wallet File Look Like?

Wallet software stores your **private key** in a local file. But if stored in plaintext, anyone who gets that file can clean you out. So it derives a key from your **passphrase** and writes it to disk with **AEAD encryption**:

```
passphrase ──KDF──► key (256 bits)
private key ────► AES-GCM ────────► salt + ciphertext   written to wallet file
```

That "salt + ciphertext" file can be stored and copied anywhere — because **without your passphrase it can't be opened**.

## The Unlocking Process (matching the code on the right)

1. read the **salt** from the file;
2. passphrase + salt → KDF → key (the same every time);
3. `${AES-GCM}` decrypts, **with authentication along the way**.

That forms a closed loop:
- **correct passphrase** → right key → decryption succeeds + auth passes → returns the private key.
- **wrong passphrase / tampered** → derived key doesn't match → AEAD authentication fails → **refused to open**.

## Why Authentication Is a Must

Suppose an attacker rewrites the ciphertext in the file (e.g. tricking you into signing with a wrong private key unknowingly) — AEAD's tag exposes it instantly. Where authentication and encryption are one, that's complete protection.

> Run the code on the right: the correct passphrase unlocks and restores the private key; a wrong passphrase makes `err` non-nil and refuses to open — that's the real protection logic behind your wallet App.