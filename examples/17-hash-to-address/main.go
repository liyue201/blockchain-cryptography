package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// b58 is the Base58 alphabet with confusable chars (0, O, I, l) removed.
const b58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// base58 encodes bytes into the Base58 alphabet.
func base58(b []byte) string {
	x := new(big.Int).SetBytes(b)
	base := big.NewInt(58)
	s := ""
	for x.Sign() > 0 {
		m := new(big.Int)
		x.DivMod(x, base, m)
		s = string(b58[m.Int64()]) + s
	}
	for _, c := range b {
		if c != 0 {
			break
		}
		s = "1" + s // a leading zero byte becomes a leading '1'
	}
	return s
}

// checksum is the first 4 bytes of double SHA-256 over the version+hash payload
// (this 4-byte tail is what makes a typo in an address detectable).
func checksum(p []byte) []byte {
	h1 := sha256.Sum256(p)
	h2 := sha256.Sum256(h1[:])
	return h2[:4]
}

func main() {
	// 1) Private key: a 256-bit random number
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Printf("1) Private key: %064x\n", priv.D)

	// 2) Public key: a point (X, Y) on the elliptic curve, compressed to 33 bytes
	pub := priv.PublicKey
	compressed := make([]byte, 33)
	compressed[0] = byte(2 + pub.Y.Bit(0)) // 02/03 prefix records the parity of Y
	pub.X.FillBytes(compressed[1:])
	fmt.Printf("2) Public key (compressed 33B): %x\n", compressed)

	// 3) Public-key hash (hash160); real chains use RIPEMD160(SHA256(pubkey))
	h := sha256.Sum256(compressed)
	hash160 := h[:20]
	fmt.Printf("3) Public-key hash (20B): %x\n", hash160)

	// 4) Build the payload the checksum is computed over: version prefix + pubkey hash
	payload := append([]byte{0x00}, hash160...)
	sum := checksum(payload)
	full := append(append([]byte{}, payload...), sum...) // version + hash + checksum
	fmt.Printf("4) payload (version+hash):     %x\n", payload)
	fmt.Printf("   checksum (first 4B):        %x\n", sum)

	// 5) Base58Check = Base58 over (version + hash + checksum)
	addr := base58(full)
	fmt.Printf("5) full (version+hash+checksum):%x\n", full)
	fmt.Printf("5) Base58Check address:         %s\n", addr)
	fmt.Println("   (Base58 strips 0/O/I/l so the address is easy to read and type, and the trailing checksum catches typos)")
}
