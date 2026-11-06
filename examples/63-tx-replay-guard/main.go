package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// Payment a transfer: nonce is the sender's "sequence number", the key to replay protection.
type Payment struct {
	From, To string
	Value    uint64
	Nonce    uint64
}

func (p Payment) hash() []byte {
	h := sha256.Sum256([]byte(p.From + p.To + itoa(p.Value) + itoa(p.Nonce)))
	return h[:]
}
func itoa(n uint64) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func main() {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := &priv.PublicKey

	// 1) a legitimate transaction: nonce=7 (this account already used 0..6)
	tx := Payment{From: "A", To: "B", Value: 100, Nonce: 7}
	dig := tx.hash()
	r, s, _ := ecdsa.Sign(rand.Reader, priv, dig)

	spent := map[uint64]bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true}
	valid := ecdsa.Verify(pub, dig, r, s) && !spent[tx.Nonce]
	fmt.Printf("1) valid signature + nonce=7 unused -> tx accepted? %v\n", valid)
	spent[tx.Nonce] = true // recorded: nonce=7 is now taken

	// 2) attacker repeats the tx verbatim (replay): nonce=7 already recorded
	used := tx
	valid2 := ecdsa.Verify(pub, used.hash(), r, s) && !spent[used.Nonce]
	fmt.Printf("2) verbatim replay of nonce=7 (already in spent) -> accepted? %v\n", valid2)

	// 3) attacker tries to roll back: replay with a smaller nonce
	old := tx
	old.Nonce = 6 // is this nonce free now? no, earlier spent[6]=true
	valid3 := ecdsa.Verify(pub, old.hash(), r, s) && !spent[old.Nonce]
	fmt.Printf("3) replay rolled back to nonce=6 (already spent) -> accepted? %v\n", valid3)

	// 4) attacker doesn't know the private key and just changes the amount -> signature fails
	evil := tx
	evil.Value = 99999
	fmt.Printf("4) change amount (no re-signing) -> signature invalid? %v\n", !ecdsa.Verify(pub, evil.hash(), r, s))

	fmt.Println("\n-> Defense-in-depth: the signature establishes 'who agreed', the nonce rejects 'duplicate / out-of-order'.")
	fmt.Println("   Missing either lets cheating through: without a signature anyone can impersonate, without nonce replay can be exploited infinitely.")
}
