package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"time"
)

type Res struct {
	Num  *big.Int
	Time time.Duration
}

var port = "127.0.0.1:8081"

func fileEncode(conn net.Conn, i int64) {

	enc := json.NewEncoder(conn)
	err := enc.Encode(i)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func fileDecode(conn net.Conn) {
	var msg Res
	dec := json.NewDecoder(conn)
	err := dec.Decode(&msg)
	if err != nil {
		fmt.Println("Decode error: ", err)
		return
	}
	fmt.Printf("%s %d\n", msg.Time, msg.Num)
}

func main() {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}

	defer conn.Close()

	fmt.Println("Set number: ")
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		i, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return
		}
		fileEncode(conn, i)
		fileDecode(conn)
	}
}
