package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// checksum uses the first 4 bytes of double SHA-256 as a checksum (prevents address typos).
func checksum(p []byte) []byte {
	h1 := sha256.Sum256(p)
	h2 := sha256.Sum256(h1[:])
	return h2[:4]
}

func main() {
	// 1) Private key: a 256-bit random number
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Printf("1) Private key: %064x\n", priv.D)

	// 2) Public key: a point (X, Y) on the elliptic curve, shown here in compressed form (33 bytes)
	pub := priv.PublicKey
	compressed := make([]byte, 33)
	compressed[0] = byte(2 + pub.Y.Bit(0)) // 02/03 prefix records the parity of Y
	pub.X.FillBytes(compressed[1:])
	fmt.Printf("2) Public key (compressed 33B): %x\n", compressed)

	// 3) Public key hash (hash160)
	//    Real chains use RIPEMD160(SHA256(pubkey)); here double SHA-256 approximates it for teaching
	h := sha256.Sum256(compressed)
	hash160 := h[:20]
	fmt.Printf("3) Public-key hash (20B): %x\n", hash160)

	// 4) Concatenate version + checksum -> payload
	payload := make([]byte, 0, 25)
	payload = append(payload, 0x00)                 // version prefix
	payload = append(payload, hash160...)           // 20-byte hash
	payload = append(payload, checksum(hash160)...) // 4-byte checksum
	fmt.Printf("4) payload (version+hash+checksum): %x\n", payload)

	fmt.Println("5) Real address = Base58Check(payload), and Base58 drops confusing chars (0/O/l/1)")
	fmt.Println("   (this example shows the core steps in hex)")
}
