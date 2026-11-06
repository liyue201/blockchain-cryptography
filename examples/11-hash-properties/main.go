package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	a := sha256.Sum256([]byte("hello blockchain"))
	b := sha256.Sum256([]byte("hello blockchin")) // we only removed the "a"

	fmt.Printf("hash of a: %x\n", a)
	fmt.Printf("hash of b: %x\n", b)

	// Property 1: Fixed length - no matter the input length, output is always 32 bytes = 256 bits
	fmt.Printf("\n[Fixed-length] both results are always %d bytes (%d bits), independent of input length\n", len(a), len(a)*8)

	// Property 2: Avalanche - a one-character change flips the output wildly
	diff := 0
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}
	fmt.Printf("[Avalanche] only 1 char changed, %d/%d bytes differ\n", diff, len(a))

	// Property 3: One-way - you can only compute forward, never invert
	fmt.Println("[One-way] knowing the output above, you cannot recover the original text")

	// Property 4: Collision resistance - finding two different inputs with the same hash is practically impossible (see lesson 8)
	fmt.Println("[Collision-resistance] finding another text that collides with a is astronomically hard in 256-bit space")
}
