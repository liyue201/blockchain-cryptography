# Ring Signatures: How Monero-Style Privacy Hides the Payer

Once an ordinary signature is public, everyone on the chain knows "who stamped it". **Ring signatures** aim to solve this: **I can prove "I'm in this group AND authorized this transaction", yet no one can tell which one I am.**

## What a "Ring" Is

The signer takes **their own public key** plus a batch of **other people's public keys** (decoys) into a set, called a **ring**. You could even cheerfully "borrow" public keys from legitimate old on-chain transactions to fill it out:

- the signer signs the message with their own private key;
- what's presented to the outside is one signature + the whole ring;
- the verifier only confirms "**some member of the ring did authorize it**", and cannot (nor should it be able to) pinpoint which one.

> Run the code on the right: the ring has 5 members, of which `[3]` is the true signer. The verify API `ringVerify` returns only a single `true` — it won't tell you the signer's index. All you see is "someone authorized it", not "who".

## Why Monero Uses It

Monero hides the sender inside a "ring", makes the recipient address a single-use **stealth address**, and locks the amount with ring confidential transactions (RingCT). So:

- the ledger is public and verifiable (anyone can inspect on-chain data);
- yet it can't be discovered **who paid whom** or **how much** — privacy comes from the "mixing" itself.

## Boundary and Truth

We must be honest: **real ring signatures have dedicated constructions** (e.g. RSA-based chained inversion, or lattice-based ring signatures) that make it **mathematically** impossible for the verifier to tell the signer from the decoys — that's the hard guarantee of "anonymity". This section, for teaching, uses ordinary ECDSA to demonstrate the conceptual skeleton of **mixing into a group of public keys, signing, and verifying without naming**; the API doesn't leak externally, which is the intuition behind anonymity, but internally trying each one still distinguishes them, so it is **not a secure implementation** — it helps you understand the *direction* of hiding, not a production-ready ring signature.