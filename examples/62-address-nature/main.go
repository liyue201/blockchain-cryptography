package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// The essence of an address = a 'fingerprint of the public key'. A person can have countless addresses, but they all come from the same public key condensed by hashing.
func main() {
	// the same key, deriving multiple "faces"
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := elliptic.Marshal(priv.Curve, priv.X, priv.Y) // uncompressed public key (prefix 04 + x + y)

	fmt.Printf("Public key (65B): %x\n", pub)
	fmt.Println()

	// fingerprint = hash of the public key (SHA-256 used here for demonstration)
	fingerprint := sha256.Sum256(pub)
	fmt.Printf("Fingerprint (32B, derives the address): %x\n", fingerprint[:])

	// another address format (version backfill, etc.) is still a fingerprint of the same key
	h := sha256.Sum256(pub)
	alt := append([]byte{0x00}, h[20:]...)
	fmt.Printf("Another address (different version/truncation): %x\n", alt[:6])

	fmt.Println("\n-> Conclusion: an address is not an account; it is a string of 'fingerprints' of a public key. On chain it acts as the external ID,")
	fmt.Println("   but what truly determines ownership is the key pair behind it (the private-key signing right = ownership).")
}
