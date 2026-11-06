package main

import (
	"crypto/sha256"
	"fmt"
)

// Block the basic unit of a blockchain: data + a hash pointing to the previous block (prevHash).
// The next block's hash contains the previous one; any change to early data 'infects' the whole tail of the chain -- the power of a hash-pointer chain.
type Block struct {
	Index    int
	Data     string
	PrevHash [32]byte
	Hash     [32]byte
}

// hashOf computes a block's hash (hashing data + prevHash together 'grabs' the previous block).
func hashOf(b Block) [32]byte {
	return sha256.Sum256([]byte(fmt.Sprintf("%d|%s|%x", b.Index, b.Data, b.PrevHash)))
}

// mineBlock makes a new block and checks that the hash really follows from the content (chain integrity).
func mineBlock(index int, data string, prev [32]byte) Block {
	b := Block{Index: index, Data: data, PrevHash: prev}
	b.Hash = hashOf(b)
	return b
}

// verifyChain checks a whole chain: each block's hash must equal its content re-hashed,
// and each block's recorded prevHash must equal the previous block's hash. Change any block and the chain breaks.
func verifyChain(blocks []Block) bool {
	for i := range blocks {
		if blocks[i].Hash != hashOf(blocks[i]) {
			return false // this block's content was touched
		}
		if i > 0 && blocks[i].PrevHash != blocks[i-1].Hash {
			return false // the link to the previous block was altered
		}
	}
	return true
}

func main() {
	// 1) build 5 blocks, each pointing at the previous block's hash
	genesis := mineBlock(0, "genesis block: reward 50", [32]byte{})
	b1 := mineBlock(1, "A pays B 5", genesis.Hash)
	b2 := mineBlock(2, "B pays C 3", b1.Hash)
	b3 := mineBlock(3, "C pays D 1", b2.Hash)
	b4 := mineBlock(4, "D pays E 0.5", b3.Hash)

	fmt.Println("Block hash chain (each block embeds the previous block's hash):")
	chain := []Block{genesis, b1, b2, b3, b4}
	for _, b := range chain {
		fmt.Printf("  #%d data=%-14s hash=%x\n       contains prevHash=%x\n", b.Index, b.Data, b.Hash[:6], b.PrevHash[:6])
	}
	fmt.Printf("Whole-chain integrity: %v\n", verifyChain(chain))

	// 2) someone secretly rewrites history in the ledger (change b1's 'A pays B 5' to 'A pays B 999')
	fake := chain
	fake[1].Data = "A pays B 999"
	fmt.Printf("After tampering with block 1, whole-chain integrity: %v (should be false: the hash pointer exposes the break all the way to the tail)\n", verifyChain(fake))

	fmt.Println("\n-> This is the 'more trustworthy than trust' ledger: change one and all change; tampering means rewriting every following block.")
}
