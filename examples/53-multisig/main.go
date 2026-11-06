package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type member struct {
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
}

// setup: pre-agree 3 members (keys). A multisig wallet = a "roster" made of several public keys.
func setup(n int) []member {
	ms := make([]member, n)
	for i := range ms {
		p, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		ms[i] = member{p, &p.PublicKey}
	}
	return ms
}

// Key point: multisig means no single member can spend the funds alone; consensus only allows spending once the threshold is reached.
func main() {
	ms := setup(3) // 3 directors
	fmt.Println("Director roster: 3 members, each with its own key pair")
	fmt.Println("Authorization rule: 2-of-3, i.e. at least 2 directors must sign for this transfer to be valid\n")

	msg := []byte("use the vault to transfer 100 ETH")
	digest := sha256.Sum256(msg)
	d := digest[:]

	// only 1 director signs
	r1, s1, _ := ecdsa.Sign(rand.Reader, ms[0].priv, d)
	sigs := 0
	if ecdsa.Verify(ms[0].pub, d, r1, s1) {
		sigs++
	}
	fmt.Printf("-> only director A signs: valid sigs %d/3 -> approved? %v (does not meet 2-of-3, rejected)\n\n",
		sigs, sigs >= 2)

	// directors A + C sign together
	rA, sA, _ := ecdsa.Sign(rand.Reader, ms[0].priv, d)
	rC, sC, _ := ecdsa.Sign(rand.Reader, ms[2].priv, d)
	nA, nC := 0, 0
	if ecdsa.Verify(ms[0].pub, d, rA, sA) {
		nA++
	}
	if ecdsa.Verify(ms[2].pub, d, rC, sC) {
		nC++
	}
	fmt.Printf("-> directors A+C sign: valid sigs %d/3 -> approved? %v (>=2, approved)\n", nA+nC, nA+nC >= 2)

	fmt.Println("\n-> 2-of-3 both prevents 'a single thief' and avoids 'single point of failure': gathering any 2 of the 3 people is enough to lawfully move funds.")
	fmt.Println("  In the real world this is the model behind multisig wallets / custodial transfers; advanced MuSig/BLS can even truly 'aggregate' the signatures into one.")
}
