package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func main() {
	// Key: AES supports 128/192/256 bits; here we use 256 bits (32 bytes)
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// GCM is the mainstream "authenticate-encrypt" mode; it is tamper-resistant by design
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// One-time nonce: every message must use a distinct random number, sent alongside the ciphertext
	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)

	msg := []byte("please encrypt this wallet private key for me")
	ct := gcm.Seal(nil, nonce, msg, nil)

	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Key (32B):    %x\n", key)
	fmt.Printf("Nonce:        %x\n", nonce)
	fmt.Printf("Ciphertext:   %x\n", ct)
	fmt.Printf("Decrypted:    %s\n", pt)
	fmt.Println("\nNote: ciphertext length = plaintext + 16-byte auth tag; you need not know the plaintext length before decrypting.")
}
