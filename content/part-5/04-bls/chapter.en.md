# Aggregate BLS Signatures: Squashing Many Signatures into One

Multisig solves "many must agree", but each additional person brings one more signature. Ethereum has tens of thousands of validators; if each stamped a signature on every block, the sync burden would be huge. **BLS compresses n signatures into 1, passing a single verification** — that's "aggregate signing".

## BLS's Core: Bilinear Pairing

Ordinary ECDSA can't do this "addition", because `(r,s)` has no way to combine multiple people's points and verify directly. BLS takes a different lane — **pairings**:

- define a pairing `e(H, a·G) = e(a·H, G)` (carrying "sums" into "products" linearly);
- **private key `sk`, public key `pk = sk·G`, message hashed into the group `H`, signature `sig = sk·H`**;
- when everyone signs the same message (e.g. a block), the signatures **add directly** into the aggregate.

> Run the code on the right: 3 signers each give a signature (3 pieces), then they combine into **1** aggregate signature, verified all at once with **1** aggregate public key. The example is "integer-based linear aggregation" — structurally of the same lineage as BLS, letting you clearly see "why it can be compressed into one", without relying on pairing math that isn't in the stdlib (the end of the text clarifies its boundary with real BLS).

## What Exactly Is Saved

- **storage**: 1000 pieces become 1, smaller blocks, faster full-node sync;
- **verification**: from verifying 1000 times to once, lighter for light clients;
- **scalability**: Bitcoin's Schnorr aggregation and Ethereum's BLS aggregation both aim for "the more participants, the less fat".

## Cost and Boundary

- pairing number theory is **brand-new math** — implementations are easier to get wrong, higher audit cost;
- careless aggregation can allow a "malicious participant to poison" (needs trusted setup / proof of possession) — engineering must take it seriously;
- this example only shows the linear-aggregation skeleton for teaching, **not a secure implementation** — real BLS needs the elliptic-curve pairing `e` and a hash-to-group map, beyond this book's hand-written scope; but the principle hides in just these words: "addition can combine, multiplication can pair".