# Spending an Address with a Signature: Transaction Authorization

The last piece of the puzzle: how exactly does a private key "spend the coins"? — **not by encrypting the transaction, but by stamping it with a signature**. Understand this, and you'll understand why the blockchain's public ledger is secure.

## What a Transaction Looks Like

The simplest transfer transaction is a few fields of plaintext:

```
from  = my address
to    = recipient address
value = amount
nonce = account sequence (anti-replay; see Part 3 §07)
```

**The transaction content is fully public and readable** — that's not a flaw but a feature: the blockchain is by design "a ledger shown to everyone".

## What Signing Does

- **Hash**: first hash the transaction content into a digest (the transaction hash);
- **Sign**: use the sender's **private key** to mathematically stamp (ECDSA) the digest, producing `(r, s)`;
- **Authorization semantics**: the signature = "I, the key holder, explicitly consent to this change";
- anyone (miners) can immediately verify with your **public key** that the stamp is yours.

> Run the code on the right and witness three cases yourself: the legitimate tx verifies `true`; **a sneakily changed amount** and **a rolled-back nonce replay** both verify `false`. Change any single byte in the content and the old stamp dies — this is how "tamper-proof + replay-proof" comes from math, not from trust.

## Why It's "Authorization", Not "Encryption"

🔐 Encryption = **hide** the content (others can't read it).

✍️ Signing = bind the content to its **author** (everyone can read it, but only you can write it).

Blockchains need the **latter**:

- the ledger is public to everyone (if encrypted, nobody could reconcile it);
- but only the private-key holder can "stamp" approval for a transfer;
- thus "public transparency + unforgeability" hold simultaneously.

## From 0 to 1: The Whole Chain a Book Can Tell

Put all four parts together, and that's everything behind one real transfer:

```
entropy → private key(41) → public key(42/3§04) → address(42)
         │
         └→ signed tx(46) ── ECDSA(3§05) / Ed25519(3§06)
                          ── hash digest (Part 1)
                          ── nonce anti-replay (3§07)
In the end it circles back: the private key is the asset; the private key is the right to sign. May you see that behind all cryptography lies a string of verifiable math.
```