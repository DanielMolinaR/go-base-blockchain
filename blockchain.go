package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// A block is a group of data
type Block struct {
	data         map[string]interface{} // Data to record on the blockchain, e.g., transaction data
	hash         string                 // the ID of the block
	previousHash string                 // The previous block’s hash, which is the cryptographic hash of the last block in the blockchain
	timestamp    time.Time              // when the block was created and added to the blockchain
	pow          int                    // amount of effort taken to derive the current block’s hash
}

// A Blockchain are multiple blocks coming together
type Blockchain struct {
	genesisBlock Block // first block added to the blockchain
	chain        []Block
	difficulty   int // minimum effort miners have to undertake to mine a block and include it in the blockchain
}

/*
Derive the block hash for the blockchain by hashing the previous block hash, current block data, timestamp, and PoW using the SHA256 algorithm.
*/
func (block Block) calculateHash() string {
	data, _ := json.Marshal(block.data)                                                                 //reads the datyafrom the block
	blockData := block.previousHash + string(data) + block.timestamp.String() + strconv.Itoa(block.pow) // Concatanate the previous block's hash, the actual block data, timestamp and Proof of work
	blockHash := sha256.Sum256([]byte(blockData))                                                       // Hash the concatenation with SHA256 algorithm
	return fmt.Sprintf("%x", blockHash)
}

/*
Mining a new block involves generating a block hash that starts with a desired number of zeros (the desired number is the mining difficulty).
This means if the difficulty of the blockchain is three, you have to generate a block hash that starts with "000" e.g., "0009a1bfb506…".
*/
func (block *Block) mine(difficulty int) {
	for !strings.HasPrefix(block.hash, strings.Repeat("0", difficulty)) {
		block.pow++
		block.hash = block.calculateHash()
	}
}

func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0", // Just one "0" since is our first Block sice there is no value for the previous hash and the data property is empty.
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

func (blockchain *Blockchain) addBlock(from, to string, amount float64) {
	// Create Block data (i.e transaction data)
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	// Retrieve the last block from the chain
	lastBlock := blockchain.chain[len(blockchain.chain)-1]

	// Create the new block
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash, // Refer the previous block variable to the last block of the chain
		timestamp:    time.Now(),
	}

	// Mine the block to create the hash based on the difficulty.
	newBlock.mine(blockchain.difficulty)
	blockchain.chain = append(blockchain.chain, newBlock)
}

// checks if the blockchain is valid so we know that no transactions have been tampered with
func (blockchain Blockchain) isValid() bool {
	for i := range blockchain.chain[1:] {
		previousBlock := blockchain.chain[i]
		currentBlock := blockchain.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

// Print blockchain's block hash
func (blockchain Blockchain) printBlocks() {
	for _, block := range blockchain.chain[1:] {
		fmt.Print(block.hash, "\n")
	}
}

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(6)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)
	blockchain.addBlock("Dani", "Bob", 27)
	blockchain.addBlock("Bob", "Dani", 50)

	// print all block hashes
	blockchain.printBlocks()

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
}
