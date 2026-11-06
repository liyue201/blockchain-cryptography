# Why Blockchain Cannot Live Without Cryptography

Imagine a public ledger: everyone writes to it, and there is no bank and no administrator. Why should you trust that:
- the ledger has not been silently altered?
- a certain sum was really sent to you by a specific person?
- nobody can impersonate you to move your money?

The answer is — **cryptography**. Blockchain, this "building of trust," is erected on just two bricks:

1. **Hashing**: turns any content into a fixed-length "fingerprint."
   Change even a single punctuation mark and the fingerprint changes beyond recognition. So "has the content been touched?" can be checked at a glance.

2. **Digital signatures**: stamp a message with a "private key," and anyone can verify the stamp using your "public key."
   So "this was indeed authorized by someone" can neither be denied nor forged.

This book is about exactly one thing: **truly understanding these two bricks — and the cryptography around them — and being able to run them with your own hands.**

> Want to try it yourself? The code on the right already computes SHA-256 for you. Try **changing a single character and clicking Run** —
> you'll see the "fingerprint" change completely. That is hashing's avalanche effect.