package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {
	alice, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &alice.PublicKey

	msg := []byte("Alice authorizes transferring out 1 BTC")
	digest := sha256.Sum256(msg)

	// 1) Alice signs with her PRIVATE key
	sig, _ := rsa.SignPKCS1v15(rand.Reader, alice, crypto.SHA256, digest[:])
	fmt.Printf("Alice's signature: %x...\n", sig[:16])

	// 2) Anyone verifies with Alice's PUBLIC key (verifying the original message -> passes)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
	fmt.Printf("verify (original message): passed? %v\n", err == nil)

	// 3) Message tampered (Eve replaced the content) -> verification fails immediately
	tampered := sha256.Sum256([]byte("Eve authorizes transferring out 1 BTC"))
	err2 := rsa.VerifyPKCS1v15(pub, crypto.SHA256, tampered[:], sig)
	fmt.Printf("verify (tampered): passed? %v  -> caught instantly\n", err2 == nil)

	fmt.Println("\n-> Signing is 'sealed by private key, verified by public key': anyone can verify, but only the key holder can sign.")
}
