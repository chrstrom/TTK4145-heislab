package localOrderDelegation

import (
	"fmt"

	"github.com/TTK4145-Students-2021/project-gruppe80/elevio"
)

type LocalOrder struct {
	floor, dir int
}

type LocalOrderDelegator struct {
	buttonInputChannel  <-chan elevio.ButtonEvent
	cabOrderChannel     chan<- int
	outsideOrderChannel chan<- LocalOrder
}

func (delegator *LocalOrderDelegator) Init(buttonCh <-chan elevio.ButtonEvent, cabOrderCh chan<- int, outsideOrderCh chan<- LocalOrder) {
	delegator.buttonInputChannel = buttonCh
	delegator.cabOrderChannel = cabOrderCh
	delegator.outsideOrderChannel = outsideOrderCh
}

func (delegator *LocalOrderDelegator) LocalOrderDelegation() {
	for {
		select {
		case buttonEvent := <-delegator.buttonInputChannel:
			if buttonEvent.Button == elevio.BT_Cab {
				//delegator.cabOrderChannel <- buttonEvent.Floor
				fmt.Println("Send cab order to floor %v", buttonEvent.Floor)
			} else {
				//order := LocalOrder{floor: buttonEvent.Floor, dir: int(buttonEvent.Button)}
				//delegator.outsideOrderChannel <- order
				fmt.Println("Send order to floor %v with direction %v", buttonEvent.Floor, buttonEvent.Button)
			}
		}
	}
}
