package main

import (
	"crypto/sha256"
	"fmt"
)

// Modeling a random oracle (RO): in the real world we use the cryptographic hash SHA-256 to play its role --
// a "deterministic but random-looking" function: the same input always gives the same output, yet the output shows no pattern.
func ro(input []byte) [32]byte { return sha256.Sum256(input) }

func commit(secret, nonce string) [32]byte { return ro([]byte(secret + "|" + nonce)) }

func main() {
	s1, n1 := "my secret", "nonce#17"
	c := commit(s1, n1)

	fmt.Println("=== Property 1: hiding ===")
	// change the nonce, the same secret yields a completely different commitment -> cannot derive the secret from c
	fmt.Printf("commit('my secret', nonce#17) = %x\n", c)
	fmt.Printf("commit('my secret', nonce#18) = %x\n", commit(s1, "nonce#18"))
	fmt.Println("-> even the same secret commits to wildly different values: the nonce keeps the commitment opaque to anyone trying to reverse it.")

	fmt.Println("\n=== Property 2: binding ===")
	fmt.Printf("Does opening with a different secret still equal the original commitment? %v\n", c == commit("his secret", n1))
	fmt.Println("-> trying to switch to another secret makes the hash mismatch -> the commitment cannot be repudiated.")

	fmt.Println("\n=== The random-oracle stamp: irregular outputs ===")
	for _, s := range []string{"a", "b", "c", "d", "e"} {
		fmt.Printf("RO('%s') = %x\n", s, ro([]byte(s)))
	}
	fmt.Println("-> a tiny input difference yields a wildly different output with no visible pattern -- this is exactly the 'random oracle' behavior hashes provide.")
}
