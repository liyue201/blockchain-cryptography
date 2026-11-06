# Post-Quantum Cryptography: Insuring Tomorrow's Assets Today

Your private-key ECDSA today leans heavily on two hard problems —**integer factorization / discrete logarithm**. But once a workable quantum computer arrives, **Shor's algorithm** could speed up both exponentially — today's "secure" signatures would go to zero in an instant. So the industry is researching **post-quantum cryptography** in advance, insuring for a decade from now while there's time.

## Two Routes

1. **Hash-based signatures**: use only one-way hashes. Hashes have no efficient quantum attack (Grover only halves the search), so they're naturally quantum-resistant. The representatives are **LMOTS / XMSS / SPHINCS+**.

> Run the code on the right — a **fully runnable Lamport one-time signature**. The core is extremely intuitive:

```
for each message bit:
  private key: reserve 2 random values  [bit=0 uses the 0th, bit=1 uses the 1st]
  public key  : hash each private-key value and publish
  signature   : for each bit of the message hash, reveal the private value for that bit
  verify      : re-hash each signature value and compare with the public key
```

> As you can see, it's `SHA-256` the whole way — no factorization, no logarithm. That's the minimal form of "post-quantum". The authentic message verifies `true`; tampered verifies `false`.

2. **Lattice-based cryptography**: based on hard problems in high-dimensional lattices (e.g. LWE/SIS). Smaller and faster than hash-based, and the main route of the new NIST standards (**ML-KEM / ML-DSA**, formerly Kyber / Dilithium).

## The Real Situation

- **harvest now, decrypt later**: attackers stockpile your ciphertext/public keys today, then mass-crack them once quantum machines mature — so "migration is coming sooner or later"; designing early is safer;
- **high migration cost**: key/signature sizes and compatibility all change hugely; the industry is paving the way gradually (compare signature sizes from Part 4 — a single Lamport signature is ~8KB).

## Boundary Note

This section uses Lamport to demonstrate the **first-principles intuition of "hash-based post-quantum"**, not a production scheme (one-time, huge keys). Real engineering pairs it with a Merkle tree to make it reusable (XMSS/SPHINCS+), or uses lattices (ML-DSA). But the target they all defend against is the same — **even if a quantum computer really arrives, your assets still have mathematics to lean on.**