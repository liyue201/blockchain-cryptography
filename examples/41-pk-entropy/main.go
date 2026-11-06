package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	// A private key = a 256-bit (32-byte) "truly random" number. Must use the cryptographically safe source crypto/rand.
	d := make([]byte, 32)
	if _, err := rand.Read(d); err != nil {
		panic(err)
	}
	n := new(big.Int).SetBytes(d)

	fmt.Printf("Private key (32 bytes / 256 bits): %x\n", d)
	fmt.Printf("Its numeric value (first 10 decimal digits): %s...\n", n.String()[:10])

	space := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	fmt.Printf("\nPrivate-key space: 2^256 ~= 10^%d\n", len(space.String())-1)
	fmt.Println("Atoms in the universe ~= 10^80, and all random bit strings ever made by humans are still far below 2^256")
	fmt.Println("-> Conclusion: the root of all cryptography is entropy. As long as the RNG is truly random, the private key cannot be brute-forced.")
	fmt.Println("\nWARNING: once the RNG goes wrong (pseudo-random/reused/leaked), no amount of space saves you - many famous thefts were random-number accidents.")
}
