package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

type Blockchain struct {
	blocks []*Block
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

func (bc *Blockchain) AddBlock(data string) *Block {
	newBlk := NewBlock(data, []byte{})
	bc.blocks = append(bc.blocks, newBlk)
	return newBlk
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func NewDb(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS blockchain(id integer PRIMARY KEY,	timestamp integer NOT NULL,	data text NOT NULL, hash text NOT NULL, prevBlockHash text NOT NULL)")
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Query("SELECT data FROM blockchain WHERE id = 1")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	if !rows.Next() {
		bc := NewGenesisBlock()
		NewItem(db, bc)
	}
}

func NewItem(db *sql.DB, bc *Block) {
	_, err := db.Exec("INSERT INTO BLOCKCHAIN (timestamp, data, hash, prevBlockHash) VALUES ($1, $2, $3, $4)",
		bc.Timestamp, bc.Data, bc.Hash, bc.PrevBlockHash)
	if err != nil {
		log.Fatalln(err)
	}
}

func lastItem(db *sql.DB) []byte {
	var id int

	rows, err := db.Query("SELECT * FROM blockchain WHERE ID=(SELECT MAX(id) FROM blockchain)")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	rows.Next()
	p := Block{}
	er := rows.Scan(&id, &p.Timestamp, &p.Data, &p.Hash, &p.PrevBlockHash)
	if er != nil {
		log.Fatalln(er)
	}
	return p.Hash
}

func ShowList(db *sql.DB) {
	var id int

	rows, err := db.Query("SELECT * FROM blockchain")
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	for rows.Next() {
		p := Block{}
		err := rows.Scan(&id, &p.Timestamp, &p.Data, &p.Hash, &p.PrevBlockHash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Prev Hash: %x\n", p.PrevBlockHash)
		fmt.Printf("Data: %s\n", p.Data)
		fmt.Printf("Hash: %x\n", p.Hash)
		fmt.Println()
	}
}

func main() {
	db, err := sql.Open("sqlite3", "blockchain.db")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	NewDb(db)
	bc := NewBlockchain()
	if len(os.Args) < 2 {
		fmt.Printf("Error: not enought arguments")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Printf("Error: not enought arguments")
		} else {
			lastHash := lastItem(db)
			bl := bc.AddBlock(os.Args[2])
			bl.PrevBlockHash = lastHash
			NewItem(db, bl)
		}
		os.Exit(0)
	case "list":
		ShowList(db)
		os.Exit(0)
	default:
		fmt.Printf("Invalid using!!!")
		os.Exit(0)
	}
}
