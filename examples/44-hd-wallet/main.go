package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"math/big"
)

var n = func() *big.Int {
	x, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	return x
}()

type extKey struct {
	priv  []byte // 32B
	chain []byte
}

// childKey: BIP32 hardened derivation. Only "private" derivation (index carries the hardened flag 0x80000000).
// data = 0x00 || ser256(parentPriv) || ser32(index); hmac-sha512(chain, data) yields IL|IR.
func childKey(parent extKey, index uint32) extKey {
	data := make([]byte, 0, 1+32+4)
	data = append(data, 0x00)
	data = append(data, parent.priv...)
	data = append(data, byte(index>>24), byte(index>>16), byte(index>>8), byte(index))
	mac := hmac.New(sha512.New, parent.chain)
	mac.Write(data)
	sum := mac.Sum(nil)
	il := new(big.Int).SetBytes(sum[:32])
	d := new(big.Int).SetBytes(parent.priv)
	childPriv := new(big.Int).Add(il, d)
	childPriv.Mod(childPriv, n)
	return extKey{priv: childPriv.FillBytes(make([]byte, 32)), chain: sum[32:]}
}

func main() {
	// A BIP32 master key normally comes from a BIP39 mnemonic via PBKDF2 to get a seed; here a random seed demonstrates the key tree.
	master := make([]byte, 32)
	rand.Read(master)
	buf := make([]byte, 1+32+32) // 0x00 || master || master (simplified as chain)
	buf[0] = 0x00
	copy(buf[1:33], master)
	copy(buf[33:], master)
	mac := hmac.New(sha512.New, master)
	mac.Write(buf[:33])
	sm := mac.Sum(nil)
	root := extKey{priv: sm[:32], chain: sm[32:]}

	fmt.Printf("Root master private key: %x\n", root.priv)
	fmt.Println("Standard path: m/44'/0'/0'/0/0   (44'=BIP44 coin, 0'=chain, 0'=account, 0=external receive, 0=address #0)")

	cur := root
	for i, idx := range []uint32{0x80000000 + 44, 0x80000000, 0x80000000, 0, 0} {
		cur = childKey(cur, idx)
		fmt.Printf("  level %d (index=%d) child priv: %x\n", i+1, idx&0x7fffffff, cur.priv)
	}
	fmt.Println("-> From one master private key, a whole tree of wallet addresses is derived; this is an HD (hierarchical deterministic) wallet.")
}
