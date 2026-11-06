# Authenticated Encryption (AEAD): Why "Decryption + Authentication" Must Be One

Encrypting alone isn't secure enough. Think about it: if an attacker **flips a few bits in your ciphertext**, you can still decrypt it, but you get **garbage** — you may not even notice you've been attacked. In large-transfers that's unacceptable.

## AEAD: Encryption That Carries a "Waybill" Along

**AEAD (Authenticated Encryption with Associated Data)** means: while decrypting, use the **authentication tag (tag)** carried with the ciphertext to check whether it has been tampered with. It's like a package that is not only sealed, but also stamped with a "seal" — opening the box reveals whether it's been pried open.

```
GCM encryption output = ciphertext + 16-byte authentication tag
   ▼ on decryption
verify the tag first → pass → decrypt
                     → fail → error immediately (not garbage)
```

## Why Can't "Verify First, Then Decrypt" Be Done Separately?

Historically many systems encrypted first then signed separately, resulting in attacks like "decryption oracle / padding oracle". Modern best practice is to **make authentication and encryption one integrated whole** — GCM and ChaCha20-Poly1305 are both AEAD — avoiding gaps to be exploited in the middle.

## What the Code Shows

- First `gcm.Open`: the ciphertext is intact, returns `err=nil`, and successfully recovers the plaintext.
- `ct[10] ^= 0x01`: only modifies 1 byte of the ciphertext.
- Second `gcm.Open`: **errors immediately**, rather than outputting garbage — that's the role of the authentication tag.

> Run the code on the right and compare the two decryptions. **Tamper with even a single bit, and AEAD catches it the moment of decryption.** That's one reason it's so often used to securely store wallet private keys.