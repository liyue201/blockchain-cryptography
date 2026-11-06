package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// mine finds a nonce such that H(header+nonce), as hex, starts with zeroes '0' characters.
func mine(header string, zeroes int) (int, string) {
	prefix := strings.Repeat("0", zeroes)
	deadline := 60_000_000
	for nonce := 0; nonce < deadline; nonce++ {
		sum := sha256.Sum256([]byte(fmt.Sprintf("%s%d", header, nonce)))
		hexs := hex.EncodeToString(sum[:])
		if strings.HasPrefix(hexs, prefix) {
			return nonce, hexs
		}
	}
	return -1, ""
}

func main() {
	for d := 1; d <= 4; d++ {
		nonce, hashOut := mine("block#1 prev=0000", d)
		if nonce < 0 {
			fmt.Printf("difficulty=%d (first 0x chars are 0, i.e. %d bits): not found before timeout\n", d, d*4)
			continue
		}
		fmt.Printf("difficulty=%d (first %d bits 0): found nonce=%d, hash=%s...\n", d, d*4, nonce, hashOut[:16])
	}
	fmt.Println("\n-> Each extra 0 (4 bits) multiplies the work by 16. That is the essence of 'proof of work': pure trial and error.")
}
