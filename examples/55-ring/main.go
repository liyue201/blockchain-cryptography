package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// A member of the ring: public key is public, private key belongs only to that person.
type member struct {
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
	hash [32]byte // public-key fingerprint, used to "show the mixing"
}

// ringVerify verifies the "whole ring" in one pass and returns only a boolean:
//
//   - msg: the raw signed message (transaction plaintext).
//   - r, s: the two big-integer components of the ECDSA signature (r, s).
//   - ring: the ring members; the signer hides among this set of public keys.
//
// It tries each member's public key in turn: as long as any single key passes ECDSA verification,
// it means "one member of this ring indeed authorized with its private key", and returns true.
// The point: the verifier only gets true/false; the interface never names which member signed,
// hiding the true signer's identity logically -- the anonymity skeleton of ring signatures ('mix in + don't name').
// NOTE: real ring signatures (e.g. Monero) are constructed so even internal step-by-step verification cannot
// distinguish the signer; this demo merely uses ordinary ECDSA to illustrate the core idea.
func ringVerify(msg []byte, r, s *big.Int, ring []member) bool {
	for _, m := range ring {
		if ok := ecdsa.Verify(m.pub, msg, r, s); ok {
			return true
		}
	}
	return false
}

func main() {
	// the ring has 5 members (Zhang, Li, Wang, Anonymous, Zhao)
	ring := make([]member, 5)
	signer := 3 // "Anonymous" will initiate the transaction
	for i := range ring {
		k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		ring[i] = member{priv: k, pub: &k.PublicKey,
			hash: sha256.Sum256(elliptic.Marshal(k.Curve, k.X, k.Y))}
	}
	fmt.Println("Ring member public-key fingerprints:")
	for i, m := range ring {
		tag := " "
		if i == signer {
			tag = "*" // only the signer himself knows it's him
		}
		fmt.Printf("  %s[%d] %x\n", tag, i, m.hash[:6])
	}

	msg := []byte("anonymously pay the recipient 10 XMR")
	// "Anonymous" signs with his private key, but presents "the whole ring"
	r, s, _ := ecdsa.Sign(rand.Reader, ring[signer].priv, msg)

	// the verifier only gets a boolean and cannot learn that signer=3
	valid := ringVerify(msg, r, s, ring)
	fmt.Printf("\nThe verifier receives 'ring + signature': is someone in the ring authorized? %v\n", valid)
	fmt.Println("But the verifier cannot see from this interface who signed (an ideal ring signature is also mathematically indistinguishable).")
}
