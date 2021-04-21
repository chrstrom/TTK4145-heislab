package localElevatorFSM

import (
	"time"

	"../cabOrderStorage"
	"../config"
	"../elevio"
	msg "../messageTypes"
)

func RunElevatorFSM(event_cabOrder <-chan int,
	fsmChannels msg.FSMChannels,
	event_floorArrival <-chan int,
	event_obstruction <-chan bool,
	event_stopButton <-chan bool,
	event_timer <-chan int) {

	elevator := initializeElevator()

	// Make sure to drive to a floor when initialized between floors
	if elevator.floor == -1 {
		elevator.direction = elevio.MD_Down
		elevio.SetMotorDirection(elevator.direction)
		elevator.state = Moving
	}

	for {
		cabOrderStorage.StoreCabOrders(elevator.requests)

		select {

		case floor := <-event_cabOrder:
			order := elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
			onRequestButtonPress(order, fsmChannels.OrderComplete, &elevator)
			setCabLights(&elevator)

		case hallOrder := <-fsmChannels.DelegateHallOrder:
			onRequestButtonPress(hallOrder, fsmChannels.OrderComplete, &elevator)
			setCabLights(&elevator)

		case costRequest := <-fsmChannels.RequestCost:

			elevatorSimulator := elevator
			cost := calculateCostForOrder(elevatorSimulator, costRequest.Order.Floor, costRequest.Order.Dir)

			if costRequest.RequestFrom == msg.Network {
				reply := costRequest.Order
				reply.Cost = cost
				fsmChannels.ReplyToNetWork <- reply
			} else {
				fsmChannels.ReplyToHallOrderManager <- cost
			}

		case newFloor := <-event_floorArrival:
			elevator.floor = newFloor

			elevio.SetFloorIndicator(newFloor)
			switch elevator.state {

			case Moving:

				if shouldStop(elevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					clearRequestAtFloor(&elevator, fsmChannels.OrderComplete)

					doorOpenTimer(&elevator)
					setCabLights(&elevator)

					elevator.state = DoorOpen
				}

			}

		case obstruction := <-event_obstruction:
			elevator.obstruction = obstruction

			if elevator.state == DoorOpen {
				elevio.SetDoorOpenLamp(true)
			}
			onDoorTimeout(&elevator)

		case <-elevator.doorTimer.C:
			onDoorTimeout(&elevator)

		}

	}
}

func initializeElevator() Elevator {
	elevator := new(Elevator)
	elevator.floor = -1
	elevator.direction = elevio.MD_Stop
	// 2D array of requests is 0 by default
	elevator.doorTimer = time.NewTimer(time.Second * config.DOOR_OPEN_DURATION)
	elevator.doorTimer.Stop()
	elevator.state = Idle
	elevator.obstruction = false

	//Load cab ordes
	cabOrders := cabOrderStorage.LoadCabOrders()
	for f := 0; f < config.N_FLOORS; f++ {
		elevator.requests[f][2] = cabOrders[f]
	}

	return *elevator
}

// Cab orders and hall orders are handled the same way by the fsm,
// but are different concepts outside of it.
func onRequestButtonPress(button_msg elevio.ButtonEvent, orderCompleteCh chan<- elevio.ButtonEvent, elevator *Elevator) {

	floor := button_msg.Floor
	button_type := button_msg.Button

	switch elevator.state {

	case DoorOpen:
		elevator.requests[floor][button_type] = true
		if elevator.floor == floor {
			clearRequestAtFloor(elevator, orderCompleteCh)
			doorOpenTimer(elevator)
		}

	case Moving:
		elevator.requests[floor][button_type] = true

	case Idle:
		if elevator.floor == floor {
			elevator.requests[floor][button_type] = true
			clearRequestAtFloor(elevator, orderCompleteCh)
			elevio.SetDoorOpenLamp(true)
			doorOpenTimer(elevator)
			elevator.state = DoorOpen
		} else {
			elevator.requests[floor][button_type] = true
			elevator.direction = chooseDirection(*elevator)
			elevio.SetMotorDirection(elevator.direction)
			elevator.state = Moving
		}
	}
}

func onDoorTimeout(elevator *Elevator) {
	if elevator.state == DoorOpen && !elevator.obstruction {
		elevio.SetDoorOpenLamp(false)
		elevator.direction = chooseDirection(*elevator)
		elevio.SetMotorDirection(elevator.direction)

		if elevator.direction == elevio.MD_Stop {
			elevator.state = Idle
		} else {
			elevator.state = Moving
		}
	}
}

func doorOpenTimer(elevator *Elevator) {
	const doorOpenTime = time.Second * config.DOOR_OPEN_DURATION
	elevio.SetDoorOpenLamp(true)
	elevator.doorTimer.Reset(doorOpenTime)
}
