package main

import (
	"fmt"
)

// BLS aggregate signatures: compress m people's signatures into one, and the chain needs just a single verification.
// Real BLS uses a bilinear pairing e(H, G) so multiple signatures can be linearly combined on an elliptic curve group.
// The standard library has no pairing; here an "integer linear aggregation" demonstrates the core structure:
//
//	group generator g is public; private key sk; public key pk = sk*g;
//	the message is first "hashed to a group element" H; signature sig = sk*H;
//	then (sum pk) and (sum sig) combine linearly, and one verification finishes the job.
type party struct{ sk int }

const g = 91 // public group generator (illustrative)

func pk(p party) int          { return p.sk * g }
func sig(p party, h int) int  { return p.sk * h }
func hashToGroup(msg int) int { return msg } // illustrative: turn a message into a group element (really a hash map)

func sum(s []int) int {
	r := 0
	for _, v := range s {
		r += v
	}
	return r
}

func main() {
	members := []party{{17}, {23}, {29}}
	msg := 12345
	H := hashToGroup(msg)

	var raw, pks []int
	for _, m := range members {
		raw = append(raw, sig(m, H))
		pks = append(pks, pk(m))
	}
	fmt.Printf("Individual signatures (3): %v\n", raw)
	fmt.Printf("Individual public keys:    %v\n", pks)

	aggSig := sum(raw) // aggregated signature = sum sig_i
	aggPk := sum(pks)  // aggregated public key = sum pk_i
	fmt.Printf("Aggregated signature (1):  %d\n", aggSig)
	fmt.Printf("Aggregated public key:     %d\n", aggPk)

	// one verification: sum pk reduced back to sum sk times H should equal the aggregated signature
	valid := (aggPk/g)*H == aggSig
	fmt.Printf("Verify once with the aggregated key: %v\n", valid)

	fmt.Println("\n-> With 1000 txs on chain: individually that's 1000 signatures; after BLS aggregation it's about 1 -- smaller blocks, faster sync.")
	fmt.Println("  Real BLS uses a bilinear pairing e to keep linearity valid on an elliptic-curve group and unforgeable; this demo shows only the structure and is never a secure implementation.")
}
