package localElevatorFSM

import (
	"time"

	"../cabOrderStorage"
	"../config"
	"../elevio"
	msg "../messageTypes"
)

// This is the driver function for the local elevator fsm
// and contains a for-select, thus should be called as a goroutine.
func RunElevatorFSM(cabOrder <-chan int,
	fsmChannels msg.FSMChannels,
	event_floorArrival <-chan int,
	event_obstruction <-chan bool,
	event_stopButton <-chan bool) {

	elevator := initializeElevator()

	for {

		cabOrderStorage.StoreCabOrders(elevator.requests)

		select {

		///////////////////////////// Order channels /////////////////////////////
		case cabOrder := <-cabOrder:
			order := elevio.ButtonEvent{Floor: cabOrder, Button: elevio.BT_Cab}
			onRequestButtonPress(order, fsmChannels.OrderComplete, &elevator)
			setCabLights(&elevator)

		case hallOrder := <-fsmChannels.DelegateHallOrder:
			onRequestButtonPress(hallOrder, fsmChannels.OrderComplete, &elevator)
			setCabLights(&elevator)

		///////////////////////////// Cost channel /////////////////////////////
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

		///////////////////////////// IO channels /////////////////////////////
		case newFloor := <-event_floorArrival:
			elevator.floor = newFloor

			elevio.SetFloorIndicator(newFloor)
			switch elevator.state {

			case Moving:
				if shouldStop(elevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					clearRequestAtFloor(&elevator, fsmChannels.OrderComplete)
					elevator.motorStopTimer.Stop()

					doorOpenTimer(&elevator)
					setCabLights(&elevator)

					elevator.state = DoorOpen
				} else {
					elevator.motorStopTimer.Reset(config.MOTOR_STOP_DETECTION_TIME)
				}

			case MotorStop:
				elevator.motorStopTimer.Stop()

				if shouldStop(elevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					clearRequestAtFloor(&elevator, fsmChannels.OrderComplete)
					elevator.motorStopTimer.Stop()

					doorOpenTimer(&elevator)
					setCabLights(&elevator)

					elevator.state = DoorOpen
				} else {
					elevator.state = Moving
					elevator.motorStopTimer.Reset(config.MOTOR_STOP_DETECTION_TIME)
				}

			}

		case obstruction := <-event_obstruction:
			elevator.obstruction = obstruction

			if elevator.state == DoorOpen {
				elevio.SetDoorOpenLamp(true)
			}

			onDoorTimeout(&elevator)

		///////////////////////////// Timeout channels /////////////////////////////
		case <-elevator.doorTimer.C:
			if elevator.obstruction {
				fsmChannels.ElevatorUnavailable <- true
				clearAllHallRequests(&elevator)
			} else {
				onDoorTimeout(&elevator)
			}

		case <-elevator.motorStopTimer.C:
			switch elevator.state {
			case Moving:
				elevator.state = MotorStop
				fsmChannels.ElevatorUnavailable <- true
				clearAllHallRequests(&elevator)
				if !requestsBelow(elevator) && !requestsAbove(elevator) {
					elevator.requests[elevator.floor+int(elevator.direction)][elevio.BT_Cab] = true
				}
				elevator.motorStopTimer.Reset(time.Second)

			case MotorStop:
				elevator.direction = chooseDirection(elevator)
				elevio.SetMotorDirection(elevator.direction)
				elevator.motorStopTimer.Reset(time.Second)
			}

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
	elevator.motorStopTimer = time.NewTimer(config.MOTOR_STOP_DETECTION_TIME)
	elevator.motorStopTimer.Stop()
	elevator.state = Idle
	elevator.obstruction = false

	//Load cab ordes
	cabOrders := cabOrderStorage.LoadCabOrders()
	for f := 0; f < config.N_FLOORS; f++ {
		elevator.requests[f][2] = cabOrders[f]
	}

	//Make sure the elevator is not between floors
	elevator.direction = elevio.MD_Down
	elevio.SetMotorDirection(elevator.direction)
	elevator.state = Moving

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

	case MotorStop:
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
			elevator.motorStopTimer.Reset(config.MOTOR_STOP_DETECTION_TIME)
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
			elevator.motorStopTimer.Reset(config.MOTOR_STOP_DETECTION_TIME)
		}
	}
}

func doorOpenTimer(elevator *Elevator) {
	const doorOpenTime = time.Second * config.DOOR_OPEN_DURATION
	elevio.SetDoorOpenLamp(true)
	elevator.doorTimer.Reset(doorOpenTime)
}
