package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

func main() {
	// Ed25519: modern choice; the signing nonce is derived from the private key + message (no extra RNG needed)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	msg := []byte("state snapshot 0x9a2f needs signing")
	sig := ed25519.Sign(priv, msg)

	fmt.Printf("Public key (32B): %x\n", pub)
	fmt.Printf("Signature (64B):  %x\n", sig)

	ok := ed25519.Verify(pub, msg, sig)
	fmt.Printf("Verify [original message]: %v\n", ok)

	// Tamper with the message -> fails
	ok2 := ed25519.Verify(pub, []byte("state snapshot 0x9999 tampered"), sig)
	fmt.Printf("Verify [tampered]: %v\n\n", ok2)

	fmt.Println("Three differences from ECDSA:")
	fmt.Println(" 1) No extra random nonce: the 'k' is derived deterministically from (private key + message), naturally preventing nonce reuse")
	fmt.Println(" 2) Constant-time: operations are independent of the data, resisting side-channel timing attacks")
	fmt.Println(" 3) Smaller signatures (64B) and faster verification -> the preferred choice for new systems")
}
