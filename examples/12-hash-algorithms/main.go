package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func main() {
	in := []byte("blockchain cryptography: from zero")

	// The same input, four algorithms each output a fingerprint of a different length/value
	md5v := md5.Sum(in)
	sha1v := sha1.Sum(in)
	sha256v := sha256.Sum256(in)
	sha512v := sha512.Sum512(in)

	fmt.Printf("%-9s %5d bits  %x\n", "MD5", md5.Size*8, md5v)
	fmt.Printf("%-9s %5d bits  %x\n", "SHA-1", sha1.Size*8, sha1v)
	fmt.Printf("%-9s %5d bits  %x\n", "SHA-256", sha256.Size*8, sha256v)
	fmt.Printf("%-9s %5d bits  %x\n", "SHA-512", sha512.Size*8, sha512v)

	fmt.Println("\nTakeaways:")
	fmt.Println("- MD5 (128-bit) and SHA-1 (160-bit) have constructible collisions -> broken, no longer usable")
	fmt.Println("- SHA-256 (256-bit) is the default for blockchain/Bitcoin, with enough security margin")
	fmt.Println("- Key question for 'is it safe': output bit length + whether a collision has been found")
}
