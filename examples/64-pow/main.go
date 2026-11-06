package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// The cryptographic essence of PoW: find a nonce such that the block's hash starts with `difficulty` zeroes.
// Finding it requires only trial and error (luck); verifying it is a single hash -- "verifiably scarce".
func mine(blockData string, difficulty int) (nonce int, hash string, tries int) {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	for {
		h := sha256.Sum256([]byte(blockData + "|" + itoa(nonce)))
		if hex.EncodeToString(h[:])[:difficulty] == prefix {
			return nonce, hex.EncodeToString(h[:]), tries
		}
		nonce++
		tries++
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func main() {
	data := "block#10: A->B 5 coins"
	for _, diff := range []int{3, 4, 5} {
		nonce, hash, tries := mine(data, diff)
		fmt.Printf("difficulty=%d  tries=%d  nonce=%d  hash=%s\n", diff, tries, nonce, hash[:16]+"...")
	}
	fmt.Println("\n-> Each +1 difficulty multiplies the expected tries by 16 -- this is where 'cost and difficulty' on-chain come from,")
	fmt.Println("   and it forces anyone rewriting history to redo the PoW for all following blocks (mining is like rowing against the current).")
}
