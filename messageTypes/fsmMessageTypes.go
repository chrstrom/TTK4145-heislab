package messageTypes

import "../elevio"

const (
	HallOrderManager = 0
	Network          = 1
)

type RequestCost struct {
	Order       OrderStamped
	RequestFrom int
}

type FSMChannels struct {
	DelegateHallOrder       chan elevio.ButtonEvent
	RequestCost             chan RequestCost
	ReplyToNetWork          chan OrderStamped
	ReplyToHallOrderManager chan int
	OrderComplete           chan elevio.ButtonEvent
	ElevatorUnavailable     chan bool
}
