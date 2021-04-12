package mock

import (
	"math/rand"
	"time"

	"../elevio"
	"../network"
)

func ReplyToRequests(request <-chan network.NewRequest, reply chan<- network.RequestReply) {
	for {
		select {
		case r := <-request:
			rep := network.RequestReply{ID: r.ID, OrderID: r.OrderID, Floor: r.Floor, Dir: r.Dir}
			rep.Cost = rand.Intn(1000)
			reply <- rep
			//fmt.Printf("	net %v - Sending mock reply \n", r.OrderID)
		}
	}
}

func ReplyToDelegations(delegation <-chan network.Delegation, reply chan<- network.DelegationConfirm) {
	for {
		select {
		case d := <-delegation:
			rep := network.DelegationConfirm{ID: d.ID, OrderID: d.OrderID, Floor: d.Floor, Dir: d.Dir}
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
