package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main() {
	// Generate a key pair (the common ECDSA pair used by wallets)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pub := &priv.PublicKey
	fmt.Println("Private key (hex):", priv.D.Text(16))

	// Sign a "transfer" message
	msg := []byte("transfer 1 BTC to Alice")
	digest := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Signature: r=%x\n           s=%x\n", r, s)

	// Verify: the original message passes
	ok := ecdsa.Verify(pub, digest[:], r, s)
	fmt.Println("Verify (original message):", ok)

	// Tampering with the message makes verification fail
	tampered := []byte("transfer 1 BTC to Eve")
	d2 := sha256.Sum256(tampered)
	fmt.Println("Verify (tampered message):", ecdsa.Verify(pub, d2[:], r, s))
}
