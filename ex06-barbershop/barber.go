package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CUTTING_TIME     = 500
	BARBERS_AMOUNT   = 10
	HALL_SITS_AMOUNT = 20
)

type Barber struct {
	val int
}

type Client struct {
	val int
}

func main() {
	clients := make(chan *Client)
	go clientProducer(clients)
	go BarberShop(clients)
	time.Sleep(2 * time.Second)
}

func clientProducer(clients chan *Client) {
	for {
		time.Sleep(time.Duration(rand.Intn(28)+7) * time.Millisecond)
		clients <- &Client{}
	}
}

func cutHear(barber *Barber, client *Client, finished chan *Barber) {
	time.Sleep(CUTTING_TIME * time.Millisecond)
	fmt.Println("Client gone")
	finished <- barber
}

func BarberShop(clients <-chan *Client) {
	freeBarbers := []*Barber{}
	waitingClient := []*Client{}
	syncBarberChan := make(chan *Barber)

	for i := 0; i < BARBERS_AMOUNT; i++ {
		freeBarbers = append(freeBarbers, &Barber{})
	}

	for {
		select {
		case client := <-clients:
			if len(freeBarbers) == 0 {
				if len(waitingClient) < HALL_SITS_AMOUNT {
					waitingClient = append(waitingClient, client)
					fmt.Printf("Client is waiting in hall (%v)\n", len(waitingClient))
				} else {
					fmt.Println("No free space for client")
				}
			} else {
				barber := freeBarbers[0]
				freeBarbers = freeBarbers[1:]
				fmt.Println("Client goes to barber")
				go cutHear(barber, client, syncBarberChan)
			}
		case barber := <-syncBarberChan:
			if len(waitingClient) > 0 {
				client := waitingClient[0]
				waitingClient = waitingClient[1:]
				fmt.Printf("Take client from room (%v)\n", len(waitingClient))
				go cutHear(barber, client, syncBarberChan)
			} else {
				fmt.Println("Barber idle")
				freeBarbers = append(freeBarbers, barber)
			}
		}
	}

}
