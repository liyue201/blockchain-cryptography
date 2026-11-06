package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Payment struct {
	To    string
	Value int
	Nonce uint64
}

func sign(priv *ecdsa.PrivateKey, p Payment) (r, s []byte) {
	h := hashOf(p)
	rr, ss, _ := ecdsa.Sign(rand.Reader, priv, h)
	return rr.Bytes(), ss.Bytes()
}

func hashOf(p Payment) []byte {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d", p.To, p.Value, p.Nonce)))
	return h[:]
}

func verify(pub *ecdsa.PublicKey, p Payment, r, s []byte) bool {
	rr, ss := new(big.Int).SetBytes(r), new(big.Int).SetBytes(s)
	return ecdsa.Verify(pub, hashOf(p), rr, ss)
}

func main() {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := &priv.PublicKey

	// Server-side in-memory "ledger/nonce pool": remember consumed nonces
	spent := map[uint64]bool{}

	commit := func(p Payment, r, s []byte) {
		if spent[p.Nonce] {
			fmt.Printf("X rejected replay: nonce=%d was already used\n", p.Nonce)
			return
		}
		if !verify(pub, p, r, s) {
			fmt.Println("X rejected: invalid signature")
			return
		}
		spent[p.Nonce] = true
		fmt.Printf("V recorded %3d -> %s (nonce=%d)\n", p.Value, p.To, p.Nonce)
	}

	// A valid payment: Carol receives 5, nonce=0
	t1 := Payment{To: "Carol", Value: 5, Nonce: 0}
	r, s := sign(priv, t1)
	commit(t1, r, s)

	// An attacker rebroadcasts the exact same transaction (replay)
	commit(t1, r, s) // expected: rejected

	// Attacker tweaks the nonce to bypass? The signature fails because the nonce changed
	evil := Payment{To: "Carol", Value: 5, Nonce: 1}
	_, _ = sign(priv, t1) // still uses the original signature
	fmt.Print("X attacker replays with changed nonce: ")
	if verify(pub, evil, r, s) {
		fmt.Println("it somehow passed")
	} else {
		fmt.Println("signature mismatch, rejected (because the nonce is part of the signature)")
	}

	// A legitimate new payment: nonce=1 increments
	t2 := Payment{To: "Carol", Value: 3, Nonce: 1}
	r2, s2 := sign(priv, t2)
	commit(t2, r2, s2) // expected: passes
}
