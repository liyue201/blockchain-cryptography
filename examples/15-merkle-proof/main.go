package main

import (
	"crypto/sha256"
	"fmt"
)

func hh(a, b [32]byte) [32]byte {
	return sha256.Sum256(append(append([]byte{}, a[:]...), b[:]...))
}
func leaf(s string) [32]byte { return sha256.Sum256([]byte(s)) }

func main() {
	txs := []string{"Tx_A", "Tx_B", "Tx_C", "Tx_D"}
	var L [4][32]byte
	for i, t := range txs {
		L[i] = leaf(t)
	}

	topL := hh(L[0], L[1])
	topR := hh(L[2], L[3])
	root := hh(topL, topR)
	fmt.Printf("The Merkle root everyone knows: %x\n\n", root[:4])

	fmt.Println("You only have: the target leaf Tx_B + one 'proof path', yet you can prove it is in the block")
	proof := [][32]byte{L[0], topR} // path: left sibling L0, right sibling topR
	cur := L[1]                     // target leaf Tx_B
	fmt.Printf("Recompute the root with it: %x\n", cur[:4])
	cur = hh(proof[0], cur) // B is the right child, first combine the left sibling
	fmt.Printf("  combine(l0, current)   -> %x\n", cur[:4])
	cur = hh(cur, proof[1]) // now we have the left child, combine the right sibling
	fmt.Printf("  combine(current, topR) -> %x\n", cur[:4])

	fmt.Printf("\nRecomputed root: %x", cur[:4])
	fmt.Printf("\nMatches the network root? -> %v\n", cur == root)
	fmt.Println("-> With only 2 hashes, we proved Tx_B really is in this tree (no need to download all transactions)")
}
