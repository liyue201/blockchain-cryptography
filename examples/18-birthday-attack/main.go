package main

import (
	"crypto/sha256"
	"fmt"
	"math"
)

// prefixKey takes the first `bits` bits (bits <= 16) of the i-th message hash.
func prefixKey(i, bits int) uint16 {
	s := sha256.Sum256([]byte(fmt.Sprintf("msg-%d", i)))
	if bits <= 8 {
		return uint16(s[0]) >> (8 - bits)
	}
	v := uint16(s[0])<<8 | uint16(s[1])
	return v >> (16 - bits)
}

func main() {
	for _, bits := range []int{8, 12, 16} {
		seen := map[uint16]int{}
		found := false
		for i := 1; i < 5_000_000; i++ {
			k := prefixKey(i, bits)
			if prev, ok := seen[k]; ok {
				// Birthday bound: expected attempts to collide on the first n bits ~= 1.253 * 2^(n/2)
				bound := 1.253 * math.Pow(2, float64(bits)/2)
				fmt.Printf("first %3d bits collide: actual %d attempts, birthday bound ~= %.0f (i=%d, %d)\n",
					bits, i, bound, prev, i)
				found = true
				break
			}
			seen[k] = i
		}
		if !found {
			fmt.Printf("first %d bits: no collision within 5M attempts\n", bits)
		}
	}
	fmt.Println("\n-> Conclusion: finding a collision needs only ~2^(n/2) tries (birthday paradox), so hashes must be long enough (e.g. 256 bits, 2^128 tries is unreachable)")
}
