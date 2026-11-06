# First Hands-On Lab: Completing a Trust Loop

Remember the opening question? **With no administrator, why should you believe that this money was really sent to you by someone?** Right now we hold exactly two tools, enough to answer it. Let's walk through the whole flow:

1. **Hash** — turn "this sentence" into a fingerprint.
2. **Sign** — stamp the fingerprint with the private key.
3. **Verify** — the whole network checks with the public key: "it was indeed authorized by the owner of that private key, and the content hasn't been touched."

## The Complete Flow (matching the code on the right)

```
Alice wants to send 1 BTC to Bob
   │  write the transaction msg
   ▼
1) hash msg ──────────► get the transaction fingerprint digest
   ▼
2) sign digest with Alice's private key ──► get the signature (r, s)
   ▼
3) broadcast to the network: msg + signature(r,s) + public key
   ▼
4) everyone verifies with the public key:
       content unchanged? √    signed by that same private key? √
       both pass → this transfer is valid
```

## Why Does "Sequence #..." Matter?

Look at step ④ in the code: change a name from `Alice` to `MALLORY` and verification immediately fails. But an experienced reader will think of this: what if someone **verbatim** rebroadcasts the same legitimate transfer ("transfer once more")? Verification would pass again! That's why in a real blockchain every transaction carries a **unique nonce/sequence number**, preventing the same transaction from being **replayed**. This anti-replay topic is covered in depth in the Part 3 signatures chapters later.

> 💡 Change the name, amount, or sequence number inside `msg` and run again, watching how the fingerprint in step ① becomes "beyond recognition" and whether verification in step ④ always fails. That's the trust assembled from "hash + signature".