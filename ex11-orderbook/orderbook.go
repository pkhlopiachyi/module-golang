package orderbook

import "sort"

type Orderbook struct {
	Ask []*Order
	Bid []*Order
}

func New() *Orderbook {
	orbook := Orderbook{}

	return &orbook
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	tr := []*Trade{}
	var or *Order
	or = nil

	switch order.Kind.String() {
	case "MARKET":
		orderbook.OrderMarket(order, &tr, &or)
	case "LIMIT":
		switch order.Side.String() {
		case "ASK":
			orderbook.OrderdAsk(order, &tr)
		case "BID":
			orderbook.OrderdBid(order, &tr)
		}
	}

	return tr, or
}

func (orderbook *Orderbook) OrderdAsk(order *Order, tr *[]*Trade) {
	flag := false

	for _, bid := range orderbook.Bid {
		if order.Price <= bid.Price {
			if order.Volume >= bid.Volume {
				newtr := Trade{Volume: bid.Volume, Price: bid.Price}
				*tr = append(*tr, &newtr)
				flag = true
				order.Volume -= bid.Volume
				bid.Volume = 0
			}
		}
	}

	if !flag {
		index := sort.Search(len(orderbook.Ask), func(i int) bool { return orderbook.Ask[i].Price > order.Price })
		orderbook.Ask = append(orderbook.Ask, order)
		copy(orderbook.Ask[index+1:], orderbook.Ask[index:])
		orderbook.Ask[index] = order

	}
}

func (orderbook *Orderbook) OrderdBid(order *Order, tr *[]*Trade) {
	flag := false

	for _, ask := range orderbook.Ask {
		if ask.Price <= order.Price {
			if ask.Volume <= order.Volume {
				newtr := Trade{Volume: ask.Volume, Price: ask.Price}
				*tr = append(*tr, &newtr)
				flag = true
				order.Volume -= ask.Volume
				ask.Volume = 0
			} else {
				newtr := Trade{Volume: order.Volume, Price: ask.Price}
				*tr = append(*tr, &newtr)
				flag = true
				ask.Volume -= order.Volume
			}
		}
	}
	if !flag {
		index := sort.Search(len(orderbook.Bid), func(i int) bool { return orderbook.Bid[i].Price < order.Price })
		orderbook.Bid = append(orderbook.Bid, order)
		copy(orderbook.Bid[index+1:], orderbook.Bid[index:])
		orderbook.Bid[index] = order
	}
}

func (orderbook *Orderbook) OrderMarket(order *Order, tr *[]*Trade, or **Order) {
	switch order.Side.String() {
	case "BID":
		for _, ask := range orderbook.Ask {
			if ask.Volume == 0 || order.Volume == 0 {
				continue
			}
			if ask.Volume <= order.Volume {
				newtr := Trade{Volume: ask.Volume, Price: ask.Price}
				*tr = append(*tr, &newtr)
				order.Volume -= ask.Volume
			} else {
				newtr := Trade{Volume: order.Volume, Price: ask.Price}
				*tr = append(*tr, &newtr)
				ask.Volume -= order.Volume
				order.Volume = 0

			}
		}
		if order.Volume > 0 {
			*or = order
		}
	case "ASK":
		newtr := Trade{}
		for _, bid := range orderbook.Bid {
			if bid.Volume == 0 || order.Volume == 0 {
				continue
			}
			if bid.Volume <= order.Volume {
				newtr = Trade{Volume: bid.Volume, Price: bid.Price}
				*tr = append(*tr, &newtr)
				order.Volume -= newtr.Volume
			} else {
				newtr = Trade{Volume: order.Volume, Price: bid.Price}
				*tr = append(*tr, &newtr)
				order.Volume = 0
			}
		}
		if order.Volume > 0 {
			*or = order
		}
	}
}
