package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// Lamport one-time signature: one of the pure hash-based post-quantum signature candidates.
// Core idea: for every bit of the message hash, pick the corresponding value from the pre-prepared private-key pair; publish the respective hashes as the public key.
// Because it relies only on the one-wayness of hashing (not integer factorization / discrete log), it is inherently resistant to quantum speed-up.

type lamport struct {
	privKey [256][2][32]byte // [bit index][bit = 0 or 1] private-key random value
	pubKey  [256][2][32]byte // hash of each private value = public key
}

func generate() *lamport {
	k := &lamport{}
	for i := 0; i < 256; i++ {
		for b := 0; b < 2; b++ {
			rand.Read(k.privKey[i][b][:])
			k.pubKey[i][b] = sha256.Sum256(k.privKey[i][b][:])
		}
	}
	return k
}

// sign: for each bit of the message hash, select the private value matching that bit -> 256 values form the signature.
func (k *lamport) sign(msg []byte) [][32]byte {
	h := sha256.Sum256(msg)
	sig := make([][32]byte, 256)
	for i := 0; i < 256; i++ {
		bit := (h[i/8] >> uint(7-i%8)) & 1
		sig[i] = k.privKey[i][bit]
	}
	return sig
}

// verify: compare each bit with the public key hash. Any mismatch -> invalid signature.
func (k *lamport) verify(msg []byte, sig [][32]byte) bool {
	h := sha256.Sum256(msg)
	for i := 0; i < 256; i++ {
		bit := (h[i/8] >> uint(7-i%8)) & 1
		if sha256.Sum256(sig[i][:]) != k.pubKey[i][bit] {
			return false
		}
	}
	return true
}

func main() {
	lp := generate()
	msg := []byte("quantum-era transfer: pay 5 ETH to B")

	sig := lp.sign(msg)
	fmt.Printf("Signature size: %d values ~= %d KB (one-time: each private key can only be used once)\n", len(sig), len(sig)*32/1000)

	fmt.Printf("Verify [original]: %v\n", lp.verify(msg, sig))
	// tampered message -> at least one bit's hash mismatches
	fmt.Printf("Verify [tampered]: %v\n", lp.verify([]byte("quantum-era transfer: pay 999 ETH to Hacker"), sig))

	fmt.Println("\n-> Why it's 'post-quantum': the signature consists only of SHA-256 (one-way hash), and hashing has no efficient quantum attack yet")
	fmt.Println("  (Shor's algorithm threatens only factorization/discrete log; Grover merely halves a 256-bit hash search).")
	fmt.Println("  Lamport is one-time and huge-keyed -> engineering couples it with a Merkle tree (MS) to make it reusable, i.e. the XMSS / SPHINCS+ family.")
	fmt.Println("  Modern NIST standards (ML-KEM/ML-DSA) add lattice crypto for smaller and faster schemes; this demo shows the industry's earliest 'hash-based' idea.")
}
