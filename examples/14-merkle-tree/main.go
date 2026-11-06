package main

import (
	"crypto/sha256"
	"fmt"
)

// treeHash hashes a single leaf.
func treeHash(d string) [32]byte { return sha256.Sum256([]byte(d)) }
func combine(a, b [32]byte) [32]byte {
	return sha256.Sum256(append(append([]byte{}, a[:]...), b[:]...))
}

// merkleRoot merges 4 leaves two by two up to the root.
func main() {
	txs := []string{"Tx_A", "Tx_B", "Tx_C", "Tx_D"}
	h := make([][32]byte, len(txs))
	for i, t := range txs {
		h[i] = treeHash(t)
		fmt.Printf("leaf %d (%s) = %x\n", i, t, h[i][:4])
	}

	// Layer 1: combine in pairs
	l1a := combine(h[0], h[1])
	l1b := combine(h[2], h[3])
	fmt.Printf("mid layer 0 (Tx_A+Tx_B) = %x\n", l1a[:4])
	fmt.Printf("mid layer 1 (Tx_C+Tx_D) = %x\n", l1b[:4])

	// Layer 2: get the Merkle root
	root := combine(l1a, l1b)
	fmt.Printf("Merkle root          = %x\n", root[:4])

	// Tampering with any transaction changes the root
	h[0] = treeHash("Tx_X")
	newRoot := combine(combine(h[0], h[1]), combine(h[2], h[3]))
	fmt.Printf("\nAfter tampering Tx_A into Tx_X, new root = %x\n", newRoot[:4])
	fmt.Printf("Did the root change? %v\n", root != newRoot)
}
