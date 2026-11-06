package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	text := "hello blockchain"
	h := sha256.Sum256([]byte(text))
	fmt.Printf("SHA-256(%q) = %s\n", text, hex.EncodeToString(h[:]))

	// Change one character and watch the hash change completely (avalanche effect)
	text2 := "hello blockchain!"
	h2 := sha256.Sum256([]byte(text2))
	fmt.Printf("SHA-256(%q) = %s\n", text2, hex.EncodeToString(h2[:]))
}
