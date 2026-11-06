package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ---- secp256k1 public domain params: y^2 = x^3 + 7 (a=0, b=7, cofactor h=1) ----
var (
	bigP = hexBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	bigN = hexBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
	GX   = hexBig("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
	GY   = hexBig("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")
)

// hexBig parses a hex string into a big integer (public params are reproducible, so they are auditable).
func hexBig(s string) *big.Int { n, _ := new(big.Int).SetString(s, 16); return n }

// pt a point; nil x means the point at infinity O.
type pt struct{ x, y *big.Int }

// addPt adds two points on the curve (handles doubling and the negated-point -> O case).
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
			return pt{} // a and b are opposites -> O
		}
		// Doubling: lambda = (3x^2 + a)/(2y), here a=0
		num := new(big.Int).Mul(a.x, a.x)
		num.Mul(num, big.NewInt(3))
		den := new(big.Int).Lsh(a.y, 1) // 2y
		lam = num.Mul(num, new(big.Int).ModInverse(den, bigP))
	} else {
		// Ordinary addition: lambda = (y2 - y1)/(x2 - x1)
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

// scalarMult computes k*G using the double-and-add scalar multiplication.
func scalarMult(k *big.Int, base pt) pt {
	var r pt // O
	t := base
	for i := k.BitLen() - 1; i >= 0; i-- {
		r = addPt(r, r)
		if k.Bit(i) == 1 {
			r = addPt(r, t)
		}
	}
	return r
}

// onCurve verifies (x,y) satisfies y^2 = x^3 + 7 (mod p).
func onCurve(p pt) bool {
	if p.x == nil {
		return true // O is always on the curve
	}
	l := new(big.Int).Mul(p.y, p.y)
	l.Mod(l, bigP)
	r := new(big.Int).Mul(p.x, p.x)
	r.Mul(r, p.x)
	r.Add(r, big.NewInt(7)).Mod(r, bigP)
	return l.Cmp(r) == 0
}

func main() {
	priv, _ := rand.Int(rand.Reader, bigN) // private key d in [1, n-1]
	pub := scalarMult(priv, pt{GX, GY})    // public key Q = d*G

	fmt.Printf("Curve secp256k1:  y^2 = x^3 + 7   (a=0, b=7, cofactor h=1)\n")
	fmt.Printf("Prime field p:    %d bits\n", bigP.BitLen())
	fmt.Printf("Base point G bits: x=%d / y=%d\n", GX.BitLen(), GY.BitLen())
	fmt.Printf("Group order n:    %d bits -> private-key space 2^%d\n", bigN.BitLen(), bigN.BitLen())
	fmt.Printf("Random private d: %064x\n", priv)
	comp := append([]byte{byte(2 + (pub.y.Bit(0) & 1))}, pub.x.Bytes()...)
	fmt.Printf("Public Q=d*G (compressed): %x\n", comp)
	fmt.Printf("Self-check: is the public key on the curve? %v\n", onCurve(pub))

	fmt.Println("\n-> Why chains choose it: minimal params (a=0, b=7), cheap computation, 256-bit security margin, and it has been heavily audited by the community.")
}
