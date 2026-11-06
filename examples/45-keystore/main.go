package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// KDF: hash the password + random salt repeatedly N times to get a fixed-length key.
// Purpose: massively raises the cost of guessing a password word-by-word - brute force must repeat this slow hash for every candidate.
func kdf(password []byte, salt []byte, iters int) []byte {
	st := append(append([]byte{}, password...), salt...)
	h := sha256.Sum256(st)
	for i := 0; i < iters; i++ {
		h = sha256.Sum256(h[:])
	}
	return h[:]
}

// seal derives a key with kdf -> AES-GCM encrypts the private key file (GCM authenticates, so a wrong password fails to decrypt).
func seal(priv, password []byte) ([]byte, []byte) {
	salt := make([]byte, 16)
	rand.Read(salt)
	key := kdf(password, salt, 200000)
	blk, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(blk)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	return append(nonce, gcm.Seal(nil, nonce, priv, nil)...), salt
}

// open: wrong password -> AES-GCM authentication fails -> explicit rejection, never a wrong return.
func open(ciphertext, salt, password []byte) ([]byte, error) {
	key := kdf(password, salt, 200000)
	blk, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(blk)
	nonce := ciphertext[:gcm.NonceSize()]
	return gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], nil)
}

func main() {
	priv := []byte("5f3a...your-private-key-hex...c0ffee") // the private key itself
	pass := []byte("correct horse battery staple")         // your wallet password

	ct, salt := seal(priv, pass)
	fmt.Printf("Encrypted wallet file (hex): %s\n", hex.EncodeToString(ct))
	fmt.Printf("Random salt:                %x\n", salt)
	fmt.Println()

	// correct password -> decrypts to the original private key
	got, err := open(ct, salt, pass)
	fmt.Printf("[correct password] unlocked? %v  private key matches? %v\n", err == nil, string(got) == string(priv))

	// wrong password -> AES-GCM fails
	if _, err2 := open(ct, salt, []byte("wrong pass")); err2 != nil {
		fmt.Printf("[wrong password] unlock failed (the desired effect): %v\n", err2)
	}
	fmt.Println("\nCold-storage idea: encrypt the wallet on an offline machine with a strong password -> only the ciphertext file goes online; online + separate key = even if servers are dumped, the private key stays locked.")
	fmt.Println("Why KDF: when a password is weak (common words), slow hashing makes every attempt cost ~a second, so brute-forcing a million passwords takes dozens of days.")
}
