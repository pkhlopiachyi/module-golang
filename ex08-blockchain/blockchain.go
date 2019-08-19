package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64  `json: "timestamp"`
	Data          []byte `json: "data"`
	PrevBlockHash []byte `json: "prevblockhash"`
	Hash          []byte `json: "hash"`
}

type Blockchain struct {
	blocks []*Block `json: "blocks"`
}

func (b *Block) SetHash() {
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

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func fCreate() {
	blkch := NewBlockchain()

	blkch.AddBlock("First")
	blkch.AddBlock("Second")
	var buf = new(bytes.Buffer)
	f, err := os.Create("hash.db.json")
	for _, block := range blkch.blocks {

		enc := json.NewEncoder(buf)
		enc.Encode(block)

		if nil != err {
			log.Fatalln(err)
		}
		defer f.Close()
		io.Copy(f, buf)
	}
}

func IsExist(fname string) {
	if _, err := os.Open(filename); err == nil {
		return false
	}

	return true
}
