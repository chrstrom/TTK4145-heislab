package orderTypes

import (
	io "../elevio"
)

type FSMChannels struct {
	DelegateHallOrder chan io.ButtonEvent
	Cost              chan int
	RequestCost       chan io.ButtonEvent
	OrderComplete     chan io.ButtonEvent
}
