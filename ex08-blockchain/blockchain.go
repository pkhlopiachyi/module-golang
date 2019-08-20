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
	Timestamp     int64  `json:"timestamp"`
	Data          []byte `json:"data"`
	PrevBlockHash []byte `json:"prevblockhash"`
	Hash          []byte `json:"hash"`
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

	block1 := NewBlockchain()

	jsonFile, err := os.Create("hash.db.json")
	if err != nil {
		fmt.Println("Message", err)
	}

	jsonWriter := io.Writer(jsonFile)
	encoder := json.NewEncoder(jsonWriter)
	err = encoder.Encode(&block1)
	if err != nil {
		fmt.Println("Message", err)
	}
}

func IsExist(fname string) bool {
	if _, err := os.Open(fname); err != nil {
		return false
	}
	return true
}
func AddItem(filename, message string) {

	var blkch Blockchain
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Message", err)
	}

	defer f.Close()
	var buf = new(bytes.Buffer)
	for _, block := range blkch.blocks {

		jsonWriter := io.Writer(f)
		enc := json.NewEncoder(jsonWriter)
		enc.Encode(block)

		if nil != err {
			log.Fatalln(err)
		}
		defer f.Close()
		io.Copy(f, buf)
	}
}

func main() {
	fileName := "hash.db.json"
	if IsExist(fileName) == false {
		fCreate()
		blkch := NewBlockchain()

		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println("Message", err)
		}
		var buf = new(bytes.Buffer)
		for _, block := range blkch.blocks {

			jsonWriter := io.Writer(f)
			enc := json.NewEncoder(jsonWriter)
			enc.Encode(block)

			if nil != err {
				log.Fatalln(err)
			}
			defer f.Close()
			io.Copy(f, buf)
		}
	}
	message := "First"
	AddItem(fileName, message)

}
