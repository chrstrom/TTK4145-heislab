package mock

import (
	"math/rand"
	"time"

	"../elevio"
	msg "../orderTypes"
)

func ReplyToRequests(request <-chan msg.OrderStamped, reply chan<- msg.OrderStamped) {
	for {
		select {
		case r := <-request:
			rep := msg.OrderStamped{ID: r.ID, OrderID: r.OrderID, Order: msg.Order{Floor: r.Order.Floor, Dir: r.Order.Dir}}
			rep.Order.Cost = rand.Intn(1000)
			reply <- rep
			//fmt.Printf("	net %v - Sending mock reply \n", r.OrderID)
		}
	}
}

func Receive(ch <-chan elevio.ButtonEvent) {
	for {
		<-ch
	}
}

func ElevatorCost(in <-chan elevio.ButtonEvent, out chan<- int) {
	for {
		<-in
		out <- rand.Intn(1000)
	}
}

func SendButtonPresses(buttons chan<- elevio.ButtonEvent, delay time.Duration) {
	e := elevio.ButtonEvent{Floor: 2, Button: elevio.BT_HallUp}
	for {
		time.Sleep(delay)
		buttons <- e
	}
}
