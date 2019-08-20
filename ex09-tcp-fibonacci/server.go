package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
)

type Res struct {
	Num  *big.Int
	Time time.Duration
}

var Server = make(map[int]*big.Int)

var port = "127.0.0.1:8081"

func Calc(n int) *big.Int {
	fn := make(map[int]*big.Int)

	for i := 0; i <= n; i++ {
		var f = big.NewInt(0)
		if i <= 2 {
			f.SetUint64(1)
		} else {
			f = f.Add(fn[i-1], fn[i-2])
		}
		fn[i] = f
	}

	return fn[n]
}

func fileEncoder(conn net.Conn, answer *Res) {
	enc := json.NewEncoder(conn)
	err := enc.Encode(answer)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func fileDecoder(conn net.Conn) (*Res, bool) {
	var num int
	dec := json.NewDecoder(conn)
	err := dec.Decode(&num)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Client number:", num)

	var res Res
	st := time.Now()
	if Server[num] != nil {
		res.Num = Server[num]
	} else {
		res.Num = Calc(num)
		Server[num] = res.Num
	}
	et := time.Since(st)
	res.Time = et

	return &res, true
}

func HandleConn(conn net.Conn) {
	fmt.Println("Accepted connecton")
	for {
		res, sucs := fileDecoder(conn)
		if !sucs {
			conn.Close()
			return
		}
		fileEncoder(conn, res)
	}
}

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Can't create server")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erroe: ", err)
			continue
		}
		go HandleConn(conn)
	}
}
