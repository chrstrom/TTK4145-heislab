package localOrderDelegation

import (
	"fmt"

	"../elevio"
)

func initializeDelegator(
	buttonCh <-chan elevio.ButtonEvent,
	hallOrderCh chan<- LocalOrder) Delegator {

	var delegator Delegator

	delegator.buttonInputChannel = buttonCh
	delegator.hallOrderChannel = hallOrderCh

	return delegator
}

func OrderDelegator(
	buttonCh <-chan elevio.ButtonEvent,
	hallOrderCh chan<- LocalOrder) {

	delegator := initializeDelegator(buttonCh, hallOrderCh)

	for {
		select {
		case buttonEvent := <-delegator.buttonInputChannel:
			if buttonEvent.Button == elevio.BT_Cab {
				//delegator.cabOrderChannel <- buttonEvent.Floor
				fmt.Println("Send cab order to floor %v", buttonEvent.Floor)
			} else {
				order := LocalOrder{Floor: buttonEvent.Floor, Dir: int(buttonEvent.Button)}
				delegator.hallOrderChannel <- order
				//fmt.Println("Send order to floor %v with direction %v", buttonEvent.Floor, buttonEvent.Button)
			}
		}
	}
}
