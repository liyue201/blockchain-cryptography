package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// ---- hand-rolled secp256k1 + ECDSA (same as Part 3, lesson 05) ----
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
		lam = new(big.Int).Mul(new(big.Int).Mul(a.x, a.x), big.NewInt(3))
		lam = lam.Mul(lam, new(big.Int).ModInverse(new(big.Int).Lsh(a.y, 1), bigP))
	} else {
		lam = new(big.Int).Sub(b.y, a.y)
		lam = lam.Mul(lam, new(big.Int).ModInverse(new(big.Int).Sub(b.x, a.x), bigP))
	}
	lam.Mod(lam, bigP)
	x3 := new(big.Int).Mul(lam, lam)
	x3.Sub(x3, a.x).Sub(x3, b.x).Mod(x3, bigP)
	y3 := new(big.Int).Mul(lam, new(big.Int).Sub(a.x, x3))
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

// signWithK signs once with a given (fixed!) k. Normal wallets generate a fresh k each time.
func signWithK(d, z, k *big.Int) (r, s *big.Int) {
	kx := scalarMult(k)
	r = new(big.Int).Mod(kx.x, bigN)
	kinv := new(big.Int).ModInverse(k, bigN)
	s = new(big.Int).Mul(r, d)
	s.Add(s, z).Mul(s, kinv).Mod(s, bigN)
	return
}

// recoverFromSameK: reproduction of ECDSA's most classic "nonce accident" --
// two messages reused the same k -> the private key can be solved directly from (r, z1, s1, z2, s2).
func recoverFromSameK(z1, s1, r, z2, s2 *big.Int) *big.Int {
	//  1) from s1-s2 == k^-1(z1-z2) recover k
	num := new(big.Int).Sub(z1, z2)
	num.Mod(num, bigN)
	den := new(big.Int).Sub(s1, s2)
	den.Mod(den, bigN)
	k := new(big.Int).Mul(num, new(big.Int).ModInverse(den, bigN))
	k.Mod(k, bigN)
	//  2) from s1*k == z1 + r*d solve d = (s1*k - z1)/r
	d := new(big.Int).Mul(s1, k)
	d.Sub(d, z1)
	d.Mul(d, new(big.Int).ModInverse(r, bigN))
	d.Mod(d, bigN)
	fmt.Printf("random nonce k recovered by the thief: ...%x\n", k.Bytes()[len(k.Bytes())-8:])
	return d
}

func main() {
	d, _ := rand.Int(rand.Reader, bigN) // the real private key (the victim)
	_ = scalarMult(d)                   // public key in play, omitted

	// attack precondition: some (bad) wallet reused the random nonce k.
	k, _ := rand.Int(rand.Reader, bigN) // generated once but reused across two transactions

	z1 := new(big.Int).SetBytes(sha256Sum("pay 1 BTC to A"))
	z2 := new(big.Int).SetBytes(sha256Sum("pay 2 BTC to B"))
	r1, s1 := signWithK(d, z1, k)
	r2, s2 := signWithK(d, z2, k)
	fmt.Printf("victim private key d: ...%x\n", d.Bytes()[len(d.Bytes())-8:])
	fmt.Printf("the two txs share the same r? %v   (r is the x-coordinate of k*G; same k -> same r -- the leak tell)\n\n", r1.Cmp(r2) == 0)

	hacked := recoverFromSameK(z1, s1, r1, z2, s2)
	fmt.Printf("private key recovered by the thief: ...%x\n", hacked.Bytes()[len(hacked.Bytes())-8:])
	fmt.Printf("matches the real private key? %v\n\n", hacked.Cmp(d) == 0)

	fmt.Println("-> Conclusion: every transaction must use its own random nonce k.")
	fmt.Println("  This is why modern implementations use RFC6979 deterministic derivation of k (Part 3, lesson 05): k is uniquely determined by the private key + message,")
	fmt.Println("  structurally eliminating 'reused/predictable nonces' -- because the most expensive mistake on-chain is a wrong nonce.")
}

func sha256Sum(s string) []byte { h := sha256.Sum256([]byte(s)); return h[:] }
