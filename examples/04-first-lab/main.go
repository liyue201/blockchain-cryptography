package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// main combines the hashing + signing from previous lessons to simulate a full
// "I received a transfer" flow.
func main() {
	// Payee Bob generates a key pair
	bobPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	bobPub := &bobPriv.PublicKey

	// Bob's bill to Alice: transfer content + a sequence number (anti-replay)
	msg := "Alice transfers 1 BTC to Bob #20261105-1"

	// 1) Hash first to get the fingerprint of this statement
	digest := sha256.Sum256([]byte(msg))
	fmt.Printf("1) Tx fingerprint (SHA-256): %x\n", digest)

	// 2) Alice signs with her private key (demo uses Bob's key as "authorizer")
	r, s, _ := ecdsa.Sign(rand.Reader, bobPriv, digest[:])
	fmt.Printf("2) Signature (r,s):\n   r=%x\n   s=%x\n", r, s)

	// 3) Anyone on the network verifies with the public key to confirm authorization
	ok := ecdsa.Verify(bobPub, digest[:], r, s)
	fmt.Printf("3) Verify with public key? %v\n", ok)

	// 4) If the message is tampered with, the signature immediately fails
	tampered := "Alice transfers 1 BTC to MALLORY #20261105-1"
	d2 := sha256.Sum256([]byte(tampered))
	fmt.Printf("4) Verify after tampering? %v\n", ecdsa.Verify(bobPub, d2[:], r, s))
}
