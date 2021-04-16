package localOrderDelegation

import "../elevio"

// This is the driver function for the local delegation node
// and contains a for-select, thus should be called as a goroutine.
func OrderDelegator(
	buttonCh <-chan elevio.ButtonEvent,
	cabOrderCh chan<- int,
	hallOrderCh chan<- LocalOrder) {

	delegator := Delegator{
		buttonInputChannel: buttonCh,
		cabOrderChannel:    cabOrderCh,
		hallOrderChannel:   hallOrderCh}

	for {
		select {
		case buttonEvent := <-delegator.buttonInputChannel:
			if buttonEvent.Button == elevio.BT_Cab {
				delegator.cabOrderChannel <- buttonEvent.Floor
			} else {
				order := LocalOrder{Floor: buttonEvent.Floor, Dir: int(buttonEvent.Button)}
				delegator.hallOrderChannel <- order
			}
		}
	}
}
