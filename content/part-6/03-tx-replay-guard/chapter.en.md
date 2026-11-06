# Transaction Signing: Ownership, Non-Repudiation & Anti-Replay

Last section we had "address = public-key fingerprint", but an address alone can't spend money. **Spending = presenting a signature that proves ownership.** These are the three most practical things in the whole book: **ownership, non-repudiation, anti-replay**.

## Signatures Solve "Ownership & Non-Repudiation"

- **Ownership**: the transaction is signed with the sender's private key; anyone verifying with their public key knows "this was authorized and sent by A", and others can't forge it;
- **Non-repudiation**: the signed-and-verified record is on-chain, so A can't later claim "I didn't send that".

## nonce Solves "Anti-Replay"

- the transaction carries a `nonce` (the next sequence number after what the sender has used);
- the chain maintains a "consumed sequence pool" per address; duplicate or out-of-order nonces are all rejected.

> Run the code on the right, four tests in a row:
> ① valid signature + unused nonce=7 → **accept**;
> ② attacker replays the same nonce=7 → **reject** (already spent);
> ③ replay with a rolled-back smaller nonce → **reject** (that sequence already consumed);
> ④ no private key, just altering the amount → **signature invalid**.
>
> Without a nonce, high-value transactions could be replayed again and again to arbitrage; without a signature, transactions could be forged as someone else's. **Signature + nonce, the two-layer guard — each alone breaks it.**

## The Anti-Fraud Combo (wrapping up what came before)

| Threat | Defense | Source |
|---|---|---|
| impersonation / repudiation | private-key signature | Part 3 |
| out-of-order / replay | transaction nonce | Part 3 §07 / here |
| tampered amount | content hash in the signature (`dig`) | Part 1 |
| stolen private key | custody / cold storage | Part 4 |

Next section adds one more "prove who did the work" layer to the whole chain — **Proof of Work (PoW)**: every block must "find a hash meeting the difficulty" to be accepted, sealing the last shortcut on the forging path.