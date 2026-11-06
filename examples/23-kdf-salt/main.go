package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// pbkdf2Like reimplements the core idea of PBKDF2 with the stdlib HMAC-SHA256:
// passphrase + salt iterated many times, slow and resistant to GPU precomputed dictionaries.
func pbkdf2Like(pass, salt string, iter, outLen int) []byte {
	out := make([]byte, 0, outLen)
	block := uint32(1)
	for len(out) < outLen {
		u := hmac.New(sha256.New, []byte(pass))
		u.Write([]byte(salt))
		u.Write([]byte{byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block)})
		t := u.Sum(nil)
		prev := append([]byte{}, t...)
		for i := 1; i < iter; i++ {
			u.Reset()
			u.Write(prev)
			prev = u.Sum(nil)
			for j := range t {
				t[j] ^= prev[j]
			}
		}
		out = append(out, t...)
		block++
	}
	return out[:outLen]
}

func main() {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)

	// Same password + same salt -> always the same 32-byte key
	k1 := pbkdf2Like("iLoveBlockchain!", string(salt), 10000, 32)
	k2 := pbkdf2Like("iLoveBlockchain!", string(salt), 10000, 32)
	fmt.Printf("salt:  %x\n", []byte(salt))
	fmt.Printf("key 1: %x\n", k1)
	fmt.Printf("key 2: %x\n", k2)
	fmt.Printf("same password + same salt -> keys equal? %v\n", string(k1) == string(k2))

	// Password off by one character / a different salt -> completely different key (unpredictable)
	k3 := pbkdf2Like("iLoveBlockchain.", string(salt), 10000, 32)
	fmt.Printf("key with last char tweaked: %x\n", k3[:8])

	// More iterations = slower -> this is the crucial "deliberate slowdown" against brute force
	fmt.Println("\n-> Iterations are the deliberate 'cost': every guessed password costs 100k computations")
}
