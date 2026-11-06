# Symmetric Encryption at a Glance: AES, One Key to Rule Them All

On-chain transactions are public, but cryptography has another half that can't do without the tool of **symmetric encryption**. It's usually used to **protect your private key, wallet files, and seed phrases off-chain**. Let's fully understand it in this section.

## What Does "Symmetric" Mean?

**Encryption and decryption use the same key.** Have the key, lock it; use the same key, unlock it.

```
plaintext ──► AES(key) ──► ciphertext
ciphertext ──► AES(key) ──► plaintext
```

- Pro: fast and secure enough; the workhorse of data encryption.
- Challenge: how to safely **hand the key to the other party**? (That's the problem "public-key cryptography" solves.)

## What Are Those Lines in the Code Doing?

1. `aes.NewCipher(key)`: instantiate AES with a 256-bit key.
2. `cipher.NewGCM(block)`: choose **GCM** mode — it not only encrypts but also attaches an **authentication tag** (which can detect tampering). This is the modern recommended usage.
3. `nonce` (one-time random number): **every message must use a different random number**, sent to the other party alongside the ciphertext; **never reuse it**, or there's a catastrophic vulnerability.
4. `gcm.Seal` encrypts / `gcm.Open` decrypts.

## One-Sentence Summary

> Symmetric encryption is like putting a padlock on a box that **the same key** can both lock and unlock. AES + GCM is currently the most reliable "lock" and is almost always seen when encrypting private keys/wallets off-chain.

> Run the code on the right and watch the key, nonce, ciphertext, and plaintext each appear. Try changing the `nonce` (so the random number after the ciphertext is wrong) and decrypting — it will immediately error.