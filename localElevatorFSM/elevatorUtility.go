package localElevatorFSM

import (
	"../config"
	"../elevio"
	types "../messageTypes"
)

func CreateFSMChannelStruct() types.FSMChannels {
	var fsmChannels types.FSMChannels

	fsmChannels.DelegateHallOrder = make(chan elevio.ButtonEvent, config.CHANNEL_BUFFER_SIZE)
	fsmChannels.ReplyToHallOrderManager = make(chan int, config.CHANNEL_BUFFER_SIZE)
	fsmChannels.ReplyToNetWork = make(chan types.OrderStamped, config.CHANNEL_BUFFER_SIZE)
	fsmChannels.RequestCost = make(chan types.RequestCost, config.CHANNEL_BUFFER_SIZE)
	fsmChannels.OrderComplete = make(chan elevio.ButtonEvent)

	return fsmChannels
}

func setCabLights(elevator *Elevator) {
	cab_button := elevio.ButtonType(2)

	for f := 0; f < config.N_FLOORS; f++ {

		if elevator.requests[f][cab_button] {
			elevio.SetButtonLamp(cab_button, f, true)
		} else {
			elevio.SetButtonLamp(cab_button, f, false)
		}

	}
}

func calculateCostForOrder(elevator Elevator, floor int, button int) int {
	var duration int = 0
	elevator.requests[floor][button] = true

	switch elevator.state {

	case Idle:
		elevator.direction = chooseDirection(elevator)
		if elevator.direction == elevio.MD_Stop {
			return duration
		}
		break
	case Moving:
		duration = config.TRAVEL_TIME / 2
		elevator.floor += int(elevator.direction)
		break
	case DoorOpen:
		duration -= config.DOOR_OPEN_DURATION / 2
	}

	for {
		if shouldStop(elevator) {
			clearRequestAtFloorSimulation(&elevator)
			if elevator.floor == floor {
				return duration
			}
			duration += config.DOOR_OPEN_DURATION
			elevator.direction = chooseDirection(elevator)

		}
		elevator.floor += int(elevator.direction)
		duration += config.TRAVEL_TIME
	}
}
