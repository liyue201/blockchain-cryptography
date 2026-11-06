package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {
	msg := []byte("hello, blockchain cryptography")

	// 1) Symmetric encryption: the same "key" both locks and unlocks
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)
	ct := gcm.Seal(nil, nonce, msg, nil)
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[1) Symmetric AES-GCM]  encrypt: %x\n   decrypt: %s\n", ct, pt)

	// 2) Public-key encryption: the public key locks, only the private key opens
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	enc, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &rsaKey.PublicKey, msg, nil)
	dec, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, enc, nil)
	fmt.Printf("[2) RSA lockbox]       locked with public key -> %s\n", dec)

	// 3) Digital signature: sign with private key, verify with public key, hash then sign
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	digest := sha256.Sum256(msg)
	r, s, _ := ecdsa.Sign(rand.Reader, priv, digest[:])
	ok := ecdsa.Verify(&priv.PublicKey, digest[:], r, s)
	fmt.Printf("[3) Signature ECDSA]   verify passed: %v\n", ok)
}
