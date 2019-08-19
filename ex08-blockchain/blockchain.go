package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     string `json: "timestamp"`
	Data          []byte `json: "data"`
	PrevBlockHash []byte `json: "prevblockhash"`
	Hash          []byte `json: "hash"`
}

type Blockchain struct {
	blocks []Block `json: "blocks"`
}

func (b Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()

	return block
}

func (blkch *Blockchain) AddBlock(data string) {
	prevBlock := blkch.blocks[len(blkch.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blkch.blocks = append(blkch.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
func main() {
	blkch := NewBlockchain()

	blkch.AddBlock("First")
	blkch.AddBlock("Second")

	for _, block := range blkch.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
	tmp := *blkch
	var buf = new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.Encode(tmp)
	f, err := os.Create("hash.db.json")
	if nil != err {
		log.Fatalln(err)
	}
	defer f.Close()
	io.Copy(f, buf)
}
