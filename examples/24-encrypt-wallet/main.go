package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// kdf simplified: passphrase + salt hashed repeatedly, mimicking the slow KDF of a real wallet.
func kdf(pass, salt string) [32]byte {
	d := sha256.Sum256([]byte(pass + "::" + salt))
	for i := 0; i < 10000; i++ {
		d = sha256.Sum256(append(d[:], []byte(pass+":"+salt)...))
	}
	return d
}

func main() {
	// The "wallet file" real content = your private key (here an example private key)
	walletData := []byte("WIF-Wallet: 5KQ3xFycWd7k2U6GdDs3r...")

	// Derive salt + key
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	saltHex := hex.EncodeToString(salt)
	key := kdf("my-strong-passphrase", saltHex)

	// Encrypt the wallet with the key using AES-GCM
	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)
	enc := gcm.Seal(nonce, nonce, walletData, nil) // nonce prepended into the file

	fmt.Printf("wallet file (salt+ciphertext): salt=%x, ciphertext=%x\n", salt, enc)

	// Next time the user unlocks with the CORRECT password
	keyOK := kdf("my-strong-passphrase", saltHex)
	b1, _ := aes.NewCipher(keyOK[:])
	g1, _ := cipher.NewGCM(b1)
	plain, err := g1.Open(nil, enc[:gcm.NonceSize()], enc[gcm.NonceSize():], nil)
	fmt.Printf("\nunlock with correct password: %s (err=%v)\n", plain, err)

	// Unlock with a WRONG password
	keyBad := kdf("wrong-passphrase-!!", saltHex)
	b2, _ := aes.NewCipher(keyBad[:])
	g2, _ := cipher.NewGCM(b2)
	_, err2 := g2.Open(nil, enc[:gcm.NonceSize()], enc[gcm.NonceSize():], nil)
	fmt.Printf("unlock with wrong password: err=%v -> authentication failed, refused\n", err2)
}
