# Digital Signatures: The "Lifeline" of Blockchain

The RSA lockbox above solves "confidentiality". But what blockchain truly can't do without is its **negative-image usage — the digital signature**. Think of it as:

> **Private key = a seal, public key = a seal-verifier.**

## The Same Keys, Just in the Opposite Direction

| Purpose | Key used | Who can do it |
|---------|----------|---------------|
| Encrypt (padlock) | recipient's public key | any sender |
| Sign (stamp) | signer's private key | only the private-key holder |
| Verify (check the seal) | signer's public key | anyone |

## Why It "Proves Three Things at Once"

Hash a message `M` with `SHA-256`, sign with the private key to get `sig`; anyone verifies with the public key. It proves, in one go:

1. **Source identity**: only the private-key holder can produce a signature that passes verification with that public key → cannot be impersonated.
2. **Untampered integrity**: change one bit of the message and verification fails.
3. **Non-repudiation**: the signature can't be pushed onto a third party — sign it and you own it.

## Why It's the Lifeline of Blockchain

Every on-chain transfer, at its core:

```
transaction content ──► SHA-256 ──► sign with private key ──► transaction carries signature
any node ──► verify with the public key for your address ──► only passes → credited
```

Without signatures, anyone could impersonate you to move assets; with them, money can only move when authorized by the private-key holder themself. **Bitcoin/Ethereum depend on this every single day.**

> Run the code on the right: in sequence see signature = private-key stamp → original message verifies → tampered message fails immediately. Replace `tampered` with any other string and the result is still "not verified".