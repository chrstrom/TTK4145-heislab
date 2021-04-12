package localOrderDelegation

import (
	"fmt"

	"../elevio"
)

type LocalOrder struct {
	Floor, Dir int
}

type LocalOrderDelegator struct {
	buttonInputChannel <-chan elevio.ButtonEvent
	cabOrderChannel    chan<- int
	hallOrderChannel   chan<- LocalOrder
}

func (delegator *LocalOrderDelegator) Init(buttonCh <-chan elevio.ButtonEvent, cabOrderCh chan<- int, hallOrderCh chan<- LocalOrder) {
	delegator.buttonInputChannel = buttonCh
	delegator.cabOrderChannel = cabOrderCh
	delegator.hallOrderChannel = hallOrderCh
}

func (delegator *LocalOrderDelegator) LocalOrderDelegation() {
	for {
		select {
		case buttonEvent := <-delegator.buttonInputChannel:
			if buttonEvent.Button == elevio.BT_Cab {
				delegator.cabOrderChannel <- buttonEvent.Floor
				fmt.Printf("Send cab order to floor %v\n", buttonEvent.Floor)
			} else {
				order := LocalOrder{Floor: buttonEvent.Floor, Dir: int(buttonEvent.Button)}
				delegator.hallOrderChannel <- order
				fmt.Printf("Send order to floor %v with direction %v\n", buttonEvent.Floor, buttonEvent.Button)
			}
		}
	}
}
