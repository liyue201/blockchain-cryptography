package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

// ---- hand-rolled secp256k1 point ops (see Part 3, lesson 04) ----
var (
	bigP = hx("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	bigN = hx("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
	G    = pt{hx("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798"),
		hx("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")}
)

func hx(s string) *big.Int { n, _ := new(big.Int).SetString(s, 16); return n }

type pt struct{ x, y *big.Int }

func addPt(a, b pt) pt {
	if a.x == nil {
		return b
	}
	if b.x == nil {
		return a
	}
	var lam *big.Int
	if a.x.Cmp(b.x) == 0 {
		if new(big.Int).Mod(new(big.Int).Add(a.y, b.y), bigP).Sign() == 0 {
			return pt{}
		}
		num := new(big.Int).Mul(a.x, a.x).Mul(new(big.Int).Mul(a.x, a.x), big.NewInt(3))
		den := new(big.Int).Lsh(a.y, 1)
		lam = num.Mul(num, new(big.Int).ModInverse(den, bigP))
	} else {
		lam = new(big.Int).Sub(b.y, a.y).Mul(new(big.Int).Sub(b.y, a.y), new(big.Int).ModInverse(new(big.Int).Sub(b.x, a.x), bigP))
	}
	lam.Mod(lam, bigP)
	x3 := new(big.Int).Mul(lam, lam)
	x3.Sub(x3, a.x).Sub(x3, b.x).Mod(x3, bigP)
	y3 := new(big.Int).Mul(new(big.Int).Sub(a.x, x3), lam)
	y3.Sub(y3, a.y).Mod(y3, bigP)
	return pt{x3, y3}
}

func scalarMult(k *big.Int) pt {
	var r pt
	t := G
	for i := k.BitLen() - 1; i >= 0; i-- {
		r = addPt(r, r)
		if k.Bit(i) == 1 {
			r = addPt(r, t)
		}
	}
	return r
}

// ---- hand-rolled Base58, the "confusable-char-stripped" encoding used for addresses ----
const b58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

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
		s = "1" + s
	}
	return s
}

// hash160 ~= RIPEMD160(SHA256(x)); here double SHA-256 approximates it for teaching
func hash160(b []byte) []byte {
	h := sha256.Sum256(b)
	h2 := sha256.Sum256(h[:])
	return h2[:20]
}

// base58checkChecksum returns the 4-byte checksum added by Base58Check:
// first 4 bytes of SHA256(SHA256(payload)), where payload = version + pubkey hash.
func base58checkChecksum(payload []byte) []byte {
	h := sha256.Sum256(payload)
	h2 := sha256.Sum256(h[:])
	return h2[:4]
}

// base58CheckValid decodes an address and returns whether its trailing checksum matches.
func base58CheckValid(s string) bool {
	lead := 0
	for _, r := range s { // leading '1's encode leading zero bytes
		if r != '1' {
			break
		}
		lead++
	}
	x := new(big.Int)
	for _, r := range s {
		m := strings.IndexRune(b58, r)
		if m < 0 {
			return false
		}
		x.Mul(x, big.NewInt(58))
		x.Add(x, big.NewInt(int64(m)))
	}
	raw := append(make([]byte, lead), x.Bytes()...)
	if len(raw) < 5 {
		return false
	}
	body, want := raw[:len(raw)-4], raw[len(raw)-4:]
	return bytes.Equal(base58checkChecksum(body), want)
}

func main() {
	d := make([]byte, 32)
	rand.Read(d)                                                     // 1) private key
	pub := scalarMult(new(big.Int).SetBytes(d))                      // 2) public key = d*G
	comp := append([]byte{byte(2 + pub.y.Bit(0))}, pub.x.Bytes()...) // compressed pubkey 33B
	h160 := hash160(comp)                                            // 3) public-key hash
	payload := append([]byte{0x00}, h160...)                         // 4) version 0x00 + pubkey hash
	checksum := base58checkChecksum(payload)                         // 4) checksum over "version+hash"
	full := append(append([]byte{}, payload...), checksum...)        // version + hash + checksum
	addr := base58(full)                                             // 5) Base58Check

	fmt.Printf("1) Private key:              %x\n", d)
	fmt.Printf("2) Public key (comp 33B):    %x\n", comp)
	fmt.Printf("3) Pubkey hash (20B):        %x\n", h160)
	fmt.Printf("4) payload = 0x00 + hash:    %x\n", payload)
	fmt.Printf("   checksum (first 4B):      %x\n", checksum)
	fmt.Printf("5) payload + checksum:       %x\n", full)
	fmt.Printf("   Base58Check address:      %s\n", addr)

	// A typo anywhere in the address is now caught by the trailing checksum.
	bad := []byte(addr)
	for i, c := range bad {
		if c != '1' {
			bad[i] = '2'
			break
		}
	}
	fmt.Printf("\nTypo test — tweak one character of %s\n", addr)
	fmt.Printf("  tampered address: %s\n", bad)
	fmt.Printf("  checksum valid?   %v\n", base58CheckValid(string(bad)))
	fmt.Println("That is the point of Base58Check: the 4-byte checksum rejects typos before they ever hit the chain.")
}
