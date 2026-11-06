package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {
	// Bob generates a key pair: public key (shareable) + private key (keep secret)
	bob, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	pub := &bob.PublicKey
	msg := []byte("transfer 0.5 BTC to Bob")

	// Anyone sending Bob a message encrypts (locks) with Bob's PUBLIC key
	ct, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, msg, nil)

	// Only Bob's PRIVATE key can decrypt (open the box)
	pt, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, bob, ct, nil)

	fmt.Printf("public key (tail): %x\n", pub.N.Bytes()[len(pub.N.Bytes())-16:])
	fmt.Printf("plaintext (unlocked by private key): %s\n", pt)
	fmt.Printf("ciphertext length: %d bytes (longer than plaintext; RSA only encrypts small messages)\n", len(ct))
	fmt.Println("\n-> Lockbox model: the public key can only 'lock', the private key only 'unlock' - the direction cannot be reversed.")
}
