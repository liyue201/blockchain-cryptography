# The Full Derivation Chain from Private Key to Address

With the private key, you can derive your entire "identity": **private key → public key → public key hash → checksum → address**. This is a **one-way**, irreversible chain — going forward anyone can verify, going backward nobody can.

## Five Steps

1. **Private key d**: a 256-bit truly random number (previous section).
2. **Public key Q = d·G**: elliptic curve scalar multiplication (learned in Part 3, §05), yielding a point on the curve. Compressed to 33 bytes (`02/03` + x-coordinate).
3. **Public key hash**: hash the public key (real chains use RIPEMD160; this tutorial approximates with double SHA-256), producing a 20-byte "fingerprint".
4. **Add version + checksum**: a `0x00` prefix marks a mainnet address; then hash the whole "prefix + hash" and take its first 4 bytes as a checksum, appended to the end. The **checksum** is what rejects an address with a single typo on-chain.
5. **Base58Check encoding**: Base58 is an alphabet that omits confusing characters (0/O, 1/l, I, etc.), making it friendlier to hand-copy than hex.

> Run the code on the right and see each step's intermediate output (private key → public key → hash → checksum → final address `17uf...`). It threads everything from Parts 1, 3, and 4 — **hashing, scalar multiplication, encoding** — into one practical chain: this is the real birth process of "addresses" on Bitcoin and many chains.

## Why It's Irreversible

The whole path uses only **one-way functions**: the discrete logarithm of scalar multiplication defies reversing, and hashing is inherently irreversible. So anyone can verify your public key/identity through your address, yet nobody can derive the private key from the address — **identity public, power private**.

## Trivia

- A checksum error means a typo is almost certainly rejected, preventing funds from being sent to a dead address.
- An address isn't an "account balance sheet" — the balance belongs to *(the private key behind the public key)* and is maintained internally by the chain's accounting.