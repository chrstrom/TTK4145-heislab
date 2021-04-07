package mock

import (
	"math/rand"
	"time"

	"../elevio"
	"../network"
)

func ReplyToRequests(request <-chan network.OrderStamped, reply chan<- network.OrderStamped) {
	for {
		select {
		case r := <-request:
			rep := network.OrderStamped{ID: r.ID, OrderID: r.OrderID, Order: network.Order{Floor: r.Order.Floor, Dir: r.Order.Dir}}
			rep.Order.Cost = rand.Intn(1000)
			reply <- rep
			//fmt.Printf("	net %v - Sending mock reply \n", r.OrderID)
		}
	}
}

func ReplyToDelegations(delegation <-chan network.OrderStamped, reply chan<- network.OrderStamped) {
	for {
		select {
		case d := <-delegation:
			rep := network.OrderStamped{ID: d.ID, OrderID: d.OrderID, Order: network.Order{Floor: d.Order.Floor, Dir: d.Order.Dir}}
			reply <- rep
			//fmt.Printf("	net %v - Sending mock confirmation \n", d.OrderID)
		}
	}
}

func SendButtonPresses(buttons chan<- elevio.ButtonEvent, delay time.Duration) {
	e := elevio.ButtonEvent{Floor: 2, Button: elevio.BT_HallUp}
	for {
		time.Sleep(delay)
		buttons <- e
	}
}
