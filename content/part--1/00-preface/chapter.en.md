# Preface: For the Reader Who Is About to "Do the Math"

This book is about **the cryptography behind blockchain**. It does not teach you how to trade coins, nor does it pile up jargon. Instead, it wants to take you on a harder, more rewarding journey: **compute, with your own hands, the mysterious math that lives on the chain.**

## Why It's Worth Reading

- **Every transfer, every block, every signature** rests on a small piece of math. Once you grasp it, you move from "having heard of blockchain" to "truly understanding how blockchain computes."
- The book is split into seven parts: hashing → symmetric encryption → digital signatures → private keys & wallets → advanced cryptography → assembling a chain. Each part is the stepping stone for the next.
- **Every concept comes with a runnable Go example** — hit "Run" on the right and watch it actually execute; you can even edit the code to test your guesses.

## Suggested Way to Read

1. **Don't panic if you're a beginner**: start with Part 0 (Getting Started) and build intuition in plain language;
2. **Don't skip the math**: each section tries to give you "a small runnable snippet of code" — running it beats memorizing ten definitions;
3. **Tinker with it**: change the hash algorithm, the signature address, or the transaction amount in the examples and run them again — understanding deepens immediately;
4. **Finally look at Part 6**: there all the pieces are assembled into a runnable trustworthy ledger — you'll see the whole book coming together.

## What This Book Can and Cannot Teach You

- **Can teach**: the cryptographic principles behind blockchain — hashing, symmetric and public-key encryption, elliptic curves and signatures, addresses and wallets, zero-knowledge, multisig/aggregate/ring signatures, and the post-quantum direction.
- **Won't cover**: consensus algorithms, UTXO, smart contracts, network synchronization, and other **system/engineering** knowledge. They are the other half of blockchain, but not the "cryptography" half.

## A One-Sentence Promise

The one conclusion this book wants you to take away lives on the final page of Part 6:

> **Security comes from verifiable math, not from trusting someone.**

Now click the first chapter in the directory on the right — let's start computing.