package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// ============ Capstone: integrate the whole book's math into a "trustworthy mini ledger" ============
//  1) hash-pointer chain (6-01)    2) Merkle to summarize txs (1-04)
//  3) private-key signing + nonce replay guard (3/6)   4) PoW hash puzzle (6-04)
// =====================================================================

// ---- accounts: the private-key signing right = asset ownership ----
var accounts = map[string]*ecdsa.PrivateKey{}

func newAccount(name string) {
	k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	accounts[name] = k
}

// ---- transactions: signature + nonce ----
type Tx struct {
	From, To string
	Value    uint64
	Nonce    uint64
	R, S     *big.Int
}

func (t Tx) digest() []byte {
	return sha256Sum([]byte(t.From + "->" + t.To + "|" + itoa(t.Value) + "|" + itoa(t.Nonce)))
}

// pay constructs a legitimate transaction signed by the from account (sender private key -> ECDSA)
func pay(from, to string, val, nonce uint64) Tx {
	t := Tx{From: from, To: to, Value: val, Nonce: nonce}
	r, s, _ := ecdsa.Sign(rand.Reader, accounts[from], t.digest())
	t.R, t.S = r, s
	return t
}

func verifyTx(t Tx) bool {
	pk := accounts[t.From]
	return ecdsa.Verify(&pk.PublicKey, t.digest(), t.R, t.S)
}

// ---- Merkle root: pair up the transaction hashes ----
func merkleRoot(txs []Tx) [32]byte {
	if len(txs) == 0 {
		return sha256.Sum256([]byte("empty"))
	}
	cur := make([][]byte, 0, len(txs))
	for i := range txs {
		cur = append(cur, txs[i].digest())
	}
	for len(cur) > 1 {
		var next [][]byte
		for i := 0; i < len(cur); i += 2 {
			right := cur[i]
			if i+1 < len(cur) {
				right = cur[i+1]
			}
			h := sha256.Sum256(append(append([]byte{}, cur[i]...), right...))
			next = append(next, h[:])
		}
		cur = next
	}
	var out [32]byte
	copy(out[:], cur[0])
	return out
}

// ---- blocks: hash-pointer chain + PoW ----
type Block struct {
	Index uint64
	Txs   []Tx
	Prev  [32]byte
	Nonce uint64
	Hash  [32]byte
}

func (b Block) body() []byte {
	return []byte(fmt.Sprintf("%d|%x|%x|%d", b.Index, merkleRoot(b.Txs), b.Prev, b.Nonce))
}
func blockHash(b Block) [32]byte { return sha256.Sum256(b.body()) }

// pack assembles a block and does PoW: find a nonce so the hash starts with `diff` zeroes
func pack(index uint64, txs []Tx, prev [32]byte, diff int) Block {
	b := Block{Index: index, Txs: txs, Prev: prev}
	for {
		h := blockHash(b)
		if hex.EncodeToString(h[:])[:diff] == zeros(diff) {
			b.Hash = h
			return b
		}
		b.Nonce++
	}
}

// verify whole-chain integrity = hash chain + PoW prefix + every signature
func verify(blocks []Block, diff int) bool {
	for i := range blocks {
		if blockHash(blocks[i]) != blocks[i].Hash {
			return false // hash was altered
		}
		if hex.EncodeToString(blocks[i].Hash[:])[:diff] != zeros(diff) {
			return false // PoW insufficient
		}
		if i > 0 && blocks[i].Prev != blocks[i-1].Hash {
			return false // chain link broken
		}
		for _, t := range blocks[i].Txs {
			if !verifyTx(t) {
				return false // a tx was not authorized by its rightful owner
			}
		}
	}
	return true
}

func itoa(n uint64) string {
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
func zeros(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "0"
	}
	return s
}
func sha256Sum(b []byte) []byte { h := sha256.Sum256(b); return h[:] }

func main() {
	newAccount("Alice")
	newAccount("Bob")
	newAccount("Carol")

	diff := 3
	// block 0: genesis (coinbase, unsigned)
	chain := []Block{pack(0, []Tx{}, [32]byte{}, diff)}

	// block 1: Alice->Bob payment (and Carol receives)
	blk1 := pack(1, []Tx{pay("Alice", "Bob", 5, 0), pay("Carol", "Alice", 2, 0)}, chain[0].Hash, diff)
	chain = append(chain, blk1)

	// block 2: two more txs (Bob's nonce needs to increase monotonically)
	blk2 := pack(2, []Tx{pay("Bob", "Carol", 1, 0)}, chain[1].Hash, diff)
	chain = append(chain, blk2)

	fmt.Println("Mini chain built:")
	for _, b := range chain {
		fmt.Printf("  block %d txs=%d PoW(nonce=%d) hash=%x\n", b.Index, len(b.Txs), b.Nonce, b.Hash[:6])
	}

	fmt.Printf("\nWhole-chain verify (hash+Merkle+signature+PoW): %v\n", verify(chain, diff))

	// attack: an attacker forges an Alice transfer -- but he does not have Alice's private key
	forged := &chain[1].Txs[0]
	forged.Value = 999999 // only change the amount, do not re-sign
	fmt.Printf("Verify after tampering with a historical amount:     %v (should be false)\n", verify(chain, diff))

	fmt.Println("\n-> The parts from five chapters are assembled here into a self-consistent trustworthy ledger.")
	fmt.Println("   Remember the sentence behind every line: security comes from 'verifiable math', not from 'trusting someone'.")
}
