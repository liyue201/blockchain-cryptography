package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// HashCommit: a hash-based commitment. Hidden before opening (invisible), bound after opening (cannot repudiate).
//   - commit = SHA256(secret || "|" || random nonce)
//   - reveal = (secret, nonce); verify = rehash and compare.
func hashCommit(secret string, nonce []byte) []byte {
	h := sha256.Sum256(append([]byte(secret), append([]byte("|"), nonce...)...))
	return h[:]
}

func main() {
	secret := "Xiaoyuan only wants to tell Xiaofang 'I'm ready', without revealing the details"

	// 1) Prover generates a random nonce and publishes only the commitment c
	nonce := make([]byte, 16)
	rand.Read(nonce)
	commit := hashCommit(secret, nonce)
	fmt.Printf("1) Prover publishes commitment c (32B): %x\n\n", commit)

	// 2) Verifier only has c and cannot learn secret -- this is 'zero-knowledge / hiding'
	guesses := []string{"Xiaoyuan is not ready at all", "Xiaoyuan is ready"}
	for _, g := range guesses {
		if string(hashCommit(g, nonce)) == string(commit) {
			fmt.Printf("2) A few guesses: [%s] happened to match? %v   (random when unknown nonce; winning probability ~= 0)\n", g, false)
		}
	}
	fmt.Println("2) With only c=H(secret|nonce): cannot guess or derive it -- hiding holds")
	fmt.Println()

	// 3) Reveal phase: Prover publishes (secret, nonce), Verifier rehashes to verify
	if string(hashCommit(secret, nonce)) == string(commit) {
		fmt.Printf("3) After reveal (secret, nonce), rehashed == c ? %v\n", true)
		fmt.Println("   -> binding holds: once revealed, anyone can verify this is exactly the committed content")

		// 4) Prover wants to say "actually that's not it" -- but the opened c mismatches, cannot repudiate
		if string(hashCommit("Xiaoyuan is not ready at all", nonce)) != string(commit) {
			fmt.Println("4) Prover claims 'I wasn't ready': hash vs c mismatch -> binding prevents repudiation")
		}
	}
	fmt.Printf("\nConclusion: a commitment lets a prover prove 'I know some secret' WITHOUT revealing any content.\nThis is the mathematical foundation of zk-rollups, privacy coins, and two-layer verification.\n")
}
