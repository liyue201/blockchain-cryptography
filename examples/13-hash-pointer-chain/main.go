package main

import (
	"crypto/sha256"
	"fmt"
)

// Block a "block": holds business data and the fingerprint (hash pointer) of the previous block.
type Block struct {
	data     string
	prevHash [32]byte // hash of the previous block
	hash     [32]byte // hash of its own content (excludes this field)
}

// fingerprint hashes a block's content (data + previous block hash).
func fingerprint(data string, prev [32]byte) [32]byte {
	buf := make([]byte, 0, 32+len(data))
	buf = append(buf, prev[:]...)
	buf = append(buf, data...)
	return sha256.Sum256(buf)
}

func main() {
	var blocks []Block
	prev := [32]byte{} // the genesis block's "prev" is agreed to be all zeros
	for _, d := range []string{"Genesis block", "Tx A: Alice->Bob", "Tx B: Bob->Carol"} {
		h := fingerprint(d, prev)
		blocks = append(blocks, Block{data: d, prevHash: prev, hash: h})
		prev = h // the next block references this block's hash
	}

	fmt.Println("--- Valid chain ---")
	for i, b := range blocks {
		fmt.Printf("block %d data=%-22s hash=%x  prev=%x\n", i, b.data, b.hash[:4], b.prevHash[:4])
	}

	// What happens if someone tampers with a middle block?
	blocks[1].data = "Tx A': Alice->EVE tampered"
	newHash := fingerprint(blocks[1].data, blocks[1].prevHash)
	fmt.Printf("\n--- After tampering with block 1 ---\n")
	fmt.Printf("block 1 new hash=%x != original hash=%x\n", newHash[:4], blocks[1].hash[:4])

	// Block 2's stored prev still points to the old hash of block 1, which no longer matches -> chain breaks
	expect := fingerprint(blocks[2].data, blocks[2].prevHash)
	fmt.Printf("block 2's own hash is still %x, and the prev it references is %x\n", expect[:4], blocks[2].prevHash[:4])
	fmt.Printf("new hash %x does not match the prev referenced by block 2 (%x) -> tampering is detected immediately\n", newHash[:4], blocks[2].prevHash[:4])
}
