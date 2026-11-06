# Anti-Replay for Signatures: How to Avoid "the Same Transaction Being Spent Twice"

Signatures solve "who authorized, and was it altered". But attackers have a sneakier trick — **replay**: they don't modify your transaction at all, they just **broadcast the same legitimate transaction verbatim again**, making the recipient/network execute it a second time. You've effectively been forced to pay twice.

## Why Is "Not Tampered" Still Dangerous?

The horror of replay is that the **signature is perfectly valid** (the content is byte-identical); no "verification" stops it. You must rely on the **transaction carrying some credential that changes**, to distinguish "I haven't seen this one before".

## The Fix: Sequence Nonce / Anti-Replay Field

Every transaction carries a **unique, monotonically increasing sequence number** (Ethereum calls it `nonce`; Bitcoin naturally resists replay via its "reference to previous spend"):

- the transaction `nonce` must **increase strictly**;
- the chain/server tracks which nonce each address "has consumed";
- receiving a duplicate or out-of-order nonce → **reject outright**, don't execute.

Key point: **the nonce must be part of the signature**. That way, even if an attacker tries "bumping the nonce up a bit to slip through", the signature fails instantly because the nonce changed — to change it they'd have to re-sign with the private key, which defeats the purpose.

> Run the code on the right to simulate a tiny "ledger":
> - first a normal tx (nonce=0) → credited;
> - the attacker **replays the identical tx verbatim** → rejected by the nonce pool;
> - the attacker **foolishly alters the nonce** to bypass → verification fails because the nonce is signed;
> - a legit new transaction (nonce=1) → credited normally.
> You'll clearly see: anti-replay relies not on "more complex signatures", but on "every transaction carrying a monotonically increasing sequence number that is itself signed".