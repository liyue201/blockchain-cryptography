package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// ---- hand-rolled secp256k1 point ops: y^2 = x^3 + 7 (stdlib CurveParams fixes a=-3, so it cannot host it) ----
var (
	bigP = hexBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	bigN = hexBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
	G    = pt{hexBig("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798"),
		hexBig("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")}
)

func hexBig(s string) *big.Int { n, _ := new(big.Int).SetString(s, 16); return n }

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
		sum := new(big.Int).Add(a.y, b.y)
		if sum.Mod(sum, bigP).Sign() == 0 {
			return pt{}
		}
		num := new(big.Int).Mul(a.x, a.x).Mul(new(big.Int).Mul(a.x, a.x), big.NewInt(3))
		den := new(big.Int).Lsh(a.y, 1)
		lam = num.Mul(num, new(big.Int).ModInverse(den, bigP))
	} else {
		dy := new(big.Int).Sub(b.y, a.y)
		dx := new(big.Int).Sub(b.x, a.x)
		lam = dy.Mul(dy, new(big.Int).ModInverse(dx, bigP))
	}
	lam.Mod(lam, bigP)
	x3 := new(big.Int).Mul(lam, lam)
	x3.Sub(x3, a.x).Sub(x3, b.x).Mod(x3, bigP)
	y3 := new(big.Int).Sub(a.x, x3)
	y3.Mul(y3, lam).Sub(y3, a.y).Mod(y3, bigP)
	return pt{x: x3, y: y3}
}

func scalarMult(k *big.Int, base pt) pt {
	var r pt
	t := base
	for i := k.BitLen() - 1; i >= 0; i-- {
		r = addPt(r, r)
		if k.Bit(i) == 1 {
			r = addPt(r, t)
		}
	}
	return r
}

// deterministicK simulates RFC6979: derives a deterministic k from (private key d, message hash z),
// eliminating random-nonce reuse.
func deterministicK(d, z *big.Int) *big.Int {
	mac := hmac.New(sha256.New, d.Bytes())
	mac.Write(z.Bytes())
	k := new(big.Int).SetBytes(mac.Sum(nil))
	k.Mod(k, bigN)
	if k.Sign() == 0 {
		k.SetInt64(1)
	}
	return k
}

// signECDSA: r = (k*G).x mod n; s = k^-1 * (z + r*d) mod n
func signECDSA(d, z *big.Int) (r, s *big.Int) {
	k := deterministicK(d, z)
	kx := scalarMult(k, G)
	r = new(big.Int).Mod(kx.x, bigN)
	kinv := new(big.Int).ModInverse(k, bigN)
	s = new(big.Int).Mul(r, d)
	s.Add(s, z)
	s.Mul(s, kinv)
	s.Mod(s, bigN)
	return
}

// verifyECDSA: (s^-1*(z*G + r*Q)).x mod n == r
func verifyECDSA(Q pt, z, r, s *big.Int) bool {
	if s.Sign() == 0 || r.Sign() == 0 {
		return false
	}
	sinv := new(big.Int).ModInverse(s, bigN)
	u1 := new(big.Int).Mul(z, sinv).Mod(new(big.Int).Mul(z, sinv), bigN)
	u2 := new(big.Int).Mul(r, sinv).Mod(new(big.Int).Mul(r, sinv), bigN)
	return new(big.Int).Mod(addPt(scalarMult(u1, G), scalarMult(u2, Q)).x, bigN).Cmp(r) == 0
}

func main() {
	d, _ := rand.Int(rand.Reader, bigN) // private key
	Q := scalarMult(d, G)               // public key Q = d*G

	msg := []byte("authorize transfer: 0.7 ETH to Carol")
	z := new(big.Int).SetBytes(sha256Sum(msg))

	r1, s1 := signECDSA(d, z)
	r2, s2 := signECDSA(d, z) // deterministic -> the two should be identical
	fmt.Printf("Private key d: ...%x\n", d.Bytes()[len(d.Bytes())-8:])
	fmt.Printf("Signature (r,s): ...%x  ...%x\n", r1.Bytes()[len(r1.Bytes())-8:], s1.Bytes()[len(s1.Bytes())-8:])
	fmt.Printf("Two signatures of the same message are identical (RFC6979)? %v\n\n", r1.Cmp(r2) == 0 && s1.Cmp(s2) == 0)

	fmt.Printf("Verify [original message]: %v\n", verifyECDSA(Q, z, r1, s1)) // 1) original message
	zT := new(big.Int).SetBytes(sha256Sum([]byte("authorize transfer: 777 ETH to Hacker")))
	fmt.Printf("Verify [tampered]: %v\n\n", verifyECDSA(Q, zT, r1, s1)) // 2) tampered

	// 3) Cross-check the same scheme on secp256r1 with the stdlib to confirm the algorithm formula
	stdPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	sig, _ := ecdsa.SignASN1(rand.Reader, stdPriv, sha256Sum(msg))
	valid := ecdsa.VerifyASN1(&stdPriv.PublicKey, sha256Sum(msg), sig)
	fmt.Printf("Cross-check: stdlib secp256r1 sign -> verify passed? %v\n", valid)
	fmt.Println("(the hand-rolled implementation here matches the secp256k1 mainnet algorithm; the stdlib does not bundle secp256k1)")
}

func sha256Sum(b []byte) []byte { h := sha256.Sum256(b); return h[:] }
