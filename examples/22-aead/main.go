package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func main() {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)

	msg := []byte("transfer: Alice gives Bob 1 BTC")
	ct := gcm.Seal(nil, nonce, msg, nil) // ciphertext automatically carries a 16-byte auth tag at the end

	// 1) Intact ciphertext: both decryption and authentication pass
	pt, err := gcm.Open(nil, nonce, ct, nil)
	fmt.Printf("[normal] decrypt: %s, authenticate: %v\n", pt, err)

	// 2) An attacker flips 1 byte in the ciphertext (even if it would not affect the decrypted value)
	ct[10] ^= 0x01
	_, err2 := gcm.Open(nil, nonce, ct, nil)
	fmt.Printf("[tampered] decryption after 1-byte flip: err = %v\n", err2)
	fmt.Println("-> The AEAD auth tag exposes any tampering the moment you decrypt, instead of outputting gibberish")
}
