package main

import (
	"crypto/rand"
	"fmt"
)

// A mini wordlist to demonstrate the BIP39 core mechanism (the real standard list has 2048 words, each giving 11 bits).
// This uses a 16-word list (4 bits per word): entropy + checksum + chunking + lookup - the same skeleton as real BIP39.
var wordlist = []string{
	"abandon", "buyer", "candy", "dolphin", "echo", "frost", "garden", "honey",
	"island", "jungle", "kite", "loyal", "mango", "novel", "oyster", "panda",
}

// wordIndex maps mnemonic words -> index for restoration
var rev = func() map[string]int {
	m := map[string]int{}
	for i, w := range wordlist {
		m[w] = i
	}
	return m
}()

// Generate a mnemonic: take 6 bytes of entropy -> checksum = (sum mod 16) -> concatenate into 13 nibbles (4 bits) -> look up words.
func toMnemonic(entropy []byte) []string {
	total := make([]byte, 0, 7)
	cs := byte(0)
	for _, b := range entropy[:6] {
		total = append(total, b)
		cs = (cs + b) % 16
	}
	h := append(total, cs)
	words := make([]string, 0, 13)
	for _, v := range h {
		hi, lo := v>>4, v&0x0f
		words = append(words, wordlist[hi], wordlist[lo])
	}
	return words
}

// Restore and verify: reconstruct entropy (first 6 bytes) from mnemonic + re-derive the checksum word to catch transcription errors.
func fromMnemonic(words []string) ([]byte, bool) {
	var entropy []byte
	var cs byte
	for i := 0; i+1 < len(words); i += 2 {
		hi, ok1 := rev[words[i]]
		lo, ok2 := rev[words[i+1]]
		if !ok1 || !ok2 {
			return nil, false
		}
		b := byte((hi << 4) | lo)
		if i < 12 { // first 6 bytes are entropy
			entropy = append(entropy, b)
			cs = (cs + b) % 16
		} else if byte(lo) != cs { // low nibble of byte 7 = checksum
			return nil, false
		}
	}
	return entropy, true
}

func main() {
	entropy := make([]byte, 6) // real BIP39 uses 16 bytes (128 bits) -> 12 words; this mini version uses 6 bytes -> 13 words
	rand.Read(entropy)
	fmt.Printf("Entropy (real 128-bit; mini demo 6B): %x\n", entropy)

	mn := toMnemonic(entropy)
	fmt.Printf("Mnemonic (%d words): %s\n", len(mn), join(mn))

	got, ok := fromMnemonic(mn)
	fmt.Printf("\nRestored entropy: %x  checksum passed? %v\n", got, ok)

	// Deliberately mistype one word -> the checksum should alarm
	bad := append([]string{}, mn...)
	bad[3] = "panda"
	_, ok2 := fromMnemonic(bad)
	fmt.Printf("Checksum after one mistyped word: %v (the wrong word is caught - this is exactly what BIP39's checksum is for)\n", ok2)
}

func join(ws []string) string {
	s := ""
	for i, w := range ws {
		if i > 0 {
			s += " "
		}
		s += w
	}
	return s
}
