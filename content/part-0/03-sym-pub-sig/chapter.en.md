# Symmetric, Public-Key, Signatures: Telling Apart the Three Tools at a Glance

In earlier chapters we met **hashing** and **signatures**. In fact the key tools that appear again and again across all of cryptography are just these three, and their differences can be summarized in one sentence:

| Tool | Plain-language idea | # of keys | Typical use |
|------|--------------------|-----------|-------------|
| **Symmetric encryption** | One key both locks and unlocks | 1 (must stay secret) | encrypting wallet files, HTTPS data transfer |
| **Public-key encryption** | The padlock is public, the key is secret | 2 (public/private) | you send me a secret letter: encrypt with my public key |
| **Digital signature** | Stamp with the private key, everyone can verify | 2 (public/private) | every on-chain transaction, proving "it's your authorization" |

## Encryption ≠ Signature (the most commonly confused point)

- **Encryption** is about "others can't see the content" — it **hides** the message.
- **Signature** is about "who did this, and was it altered" — it **stamps** the message.

On-chain transactions are **not encrypted** (the chain is public); they use **signatures**, whose purpose is authenticity and non-repudiation.

## What the Three Tools in the Code Do

1. **AES symmetric encryption**: the same `key` encrypts with `gcm.Seal` and decrypts with `gcm.Open`.
2. **RSA public-key encryption**: encrypt with `rsaKey.PublicKey` (public-key padlock), only `rsaKey` (the private key) can open it.
3. **ECDSA signature**: hash the message, sign with the private key to get `(r, s)`, and anyone verifies with the public key via `Verify`.

> 💡 See the output in the top-right? Each of the three lines is one of the three tools. Try changing `msg` and running again —
> the symmetric encryption and the signature outputs all change, whereas RSA encryption decrypts back to the same plaintext.