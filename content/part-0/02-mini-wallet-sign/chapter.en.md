# Hashing: The Content Fingerprint

> Before we move on, let's use a small experiment closer to a "wallet" to tie two concepts together:
> generate a key pair → sign with the private key → verify with the public key.

## How to Read This Code

1. **Generate the key pair**: `ecdsa.GenerateKey` produces at once a "private key" (secret) and a "public key" (public).
2. **Sign**: we hash a "transfer message", then sign it with the private key to get `(r, s)`.
3. **Verify**: anyone holding the public key can check whether the signature was made by the matching private key.

When you run it, you'll see two key conclusions:

- **Original message** → `verified: true`: this was indeed authorized by the private-key holder.
- **Tampered message** → `verified: false`: the moment the content changes, the signature breaks.

## That's Why On-Chain Records Need Signatures

On the chain, every transfer carries the sender's digital signature. **Without the private key nobody can impersonate you to transfer funds** — that guarantee is provided by signatures.

> 💡 Try changing the name inside `msg` (for example, turn `Alice` into `Bob`),
> or tweak the content of the `tampered` line, then click "Run" and watch how the verification result changes.