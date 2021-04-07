package localOrderDelegation

import 	"../elevio"

type LocalOrder struct {
	Floor, Dir int
}

type Delegator struct {
	buttonInputChannel <-chan elevio.ButtonEvent
	cabOrderChannel    chan<- int
	hallOrderChannel   chan<- LocalOrder
}