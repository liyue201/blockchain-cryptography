package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// A mini transaction: its meaning is "I (private key) authorize moving value coins from from to to"; nonce prevents replay.
type Tx struct {
	From, To string
	Value    uint64
	Nonce    uint64
}

func (t Tx) hash() []byte {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d|%d", t.From, t.To, t.Value, t.Nonce)))
	return h[:]
}

// Key insight: signing is not "wrapping" the transaction; it is a mathematical stamp applied with the private key to the transaction's hash.
// Anyone can read the plaintext (from/to/value/nonce), but without the private key nobody can forge a valid stamp.
func main() {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := &priv.PublicKey

	// a legitimate transaction initiated by the user
	orig := Tx{From: "address-A", To: "address-B", Value: 0.5e18, Nonce: 7}
	digest := orig.hash()
	r, s, _ := ecdsa.Sign(rand.Reader, priv, digest)
	fmt.Printf("Transaction plaintext: from=%s -> to=%s  value=%d  nonce=%d\n",
		orig.From, orig.To, orig.Value, orig.Nonce)
	fmt.Printf("Miner verifies with address-A's public key: %v\n\n", ecdsa.Verify(pub, digest, r, s))

	// want to silently change the amount / roll back the nonce and replay?
	tampered := orig
	tampered.Value = 1e18 // change 0.5 to 1
	d2 := tampered.hash()
	fmt.Printf("Verify tampered-value tx: %v  (content changed, the old stamp must fail)\n", ecdsa.Verify(pub, d2, r, s))

	replayed := orig
	replayed.Nonce = 6 // replay an earlier-sequence tx
	d3 := replayed.hash()
	fmt.Printf("Verify rolled-back-nonce tx: %v  (also fails, preventing out-of-order replay)\n", ecdsa.Verify(pub, d3, r, s))

	fmt.Println("\n-> Conclusion: a signature is 'authorization', not 'encryption'. Miners recognize only the private key's signature; the contents can be public - this is the cornerstone that lets a public ledger operate safely.")
	_ = pub
}
