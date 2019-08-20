package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Block struct {
	Timestamp     string `json:"timestamp"`
	Data          string `json:"data"`
	Hash          string `json:"hash"`
	PrevBlockHash string `json:"prevblockhash"`
}

type Blockchain struct {
	Blocks []Block `json:"blocks"`
}

func newBlockchain() *Blockchain {
	return &Blockchain{[]Block{genesisBlock()}}
}

func genesisBlock() Block {
	t := time.Now()
	block := Block{t.String(), "Genesis Block", "", ""}
	return block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Hash, data)
	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Println("")
	fmt.Println("Hash generated:  ", newBlock.Hash)
	fmt.Println("")
}

func NewBlock(prevBlockHash string, data string) Block {
	t := time.Now()
	block := Block{t.String(), data, "", prevBlockHash}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	record := b.Timestamp + b.PrevBlockHash + b.Data
	sha := sha256.New()
	sha.Write([]byte(record))
	hashed := sha256.Sum256(nil)
	b.Hash = hex.EncodeToString(hashed)
}

func ReadFile(fileName string) Blockchain {
	var blockchain Blockchain
	buffer, _ := ioutil.ReadFile("blockchain.json")
	json.Unmarshal(buffer, &blockchain)
	return blockchain
}

func WriteFile(filename string, data string) {
	_, err := os.Open(filename)
	if err != nil {
		fmt.Print("Creating blockchain\n")
		CreateFile()
	}
	blockchain := ReadFile(filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	Must(err)
	writer := io.Writer(file)
	encoder := json.NewEncoder(writer)
	blockchain.AddBlock(data)
	err = encoder.Encode(&blockchain)
	Must(err)
}

func CreateFile() {
	blockchain := newBlockchain()
	file, err := os.Create("blockchain.json")
	Must(err)
	fmt.Print("Blockchain created\n")
	writer := io.Writer(file)
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(&blockchain)
	Must(err)
}

func Must(err error) {
	if err != nil {
		fmt.Println("Something wrong: ", err)
	}
}

func List() {
	blockchain := ReadFile("blockchain.json")
	for i := range blockchain.Blocks {
		fmt.Println("Time: ", blockchain.Blocks[i].Timestamp)
		fmt.Println("Data: ", blockchain.Blocks[i].Data)
		fmt.Println("Previous Hash: ", blockchain.Blocks[i].PrevBlockHash)
		fmt.Println("")
		fmt.Println("Hash: ", blockchain.Blocks[i].Hash)
		fmt.Println("")
		fmt.Println("")
	}
}

func main() {
	switch os.Args[1] {
	case "add":
		WriteFile("blockchain.json", os.Args[2])
	case "list":
		List()
	default:
		fmt.Println("")
		fmt.Println("Wrong argument")
		fmt.Println("")
	}
}
