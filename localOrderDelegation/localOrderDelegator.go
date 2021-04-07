package localOrderDelegation

import (
	"fmt"

	"../elevio"
)

func initializeDelegator(
	buttonCh <-chan elevio.ButtonEvent,
	hallOrderCh chan<- LocalOrder,
	cabOrderCh chan<- int) Delegator {

	var delegator Delegator

	delegator.buttonInputChannel = buttonCh
	delegator.hallOrderChannel = hallOrderCh
	delegator.cabOrderChannel = cabOrderCh

	return delegator
}

func OrderDelegator(
	buttonCh <-chan elevio.ButtonEvent,
	hallOrderCh chan<- LocalOrder,
	cabOrderCh chan<- int) {

	delegator := initializeDelegator(buttonCh, hallOrderCh, cabOrderCh)

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
