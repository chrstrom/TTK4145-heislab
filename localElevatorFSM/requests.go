package localElevatorFSM

import (
	"../config"
	"../elevio"
)

func requestsAbove(e Elevator) bool {
	for floor := e.floor + 1; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if e.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for floor := 0; floor < e.floor; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			if e.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func chooseDirection(e Elevator) elevio.MotorDirection {
	switch e.direction {

	// Note which functions are called first for each case!
	// For MD_Stop it doesn't really matter which one goes first
	case elevio.MD_Up:
		if requestsAbove(e) {
			return elevio.MD_Up
		} else if requestsBelow(e) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Stop:

		fallthrough

	case elevio.MD_Down:
		if requestsBelow(e) {
			return elevio.MD_Down
		} else if requestsAbove(e) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}

	}
	return elevio.MD_Stop
}

func shouldStop(e Elevator) bool {

	var floor = e.floor

	switch e.direction {

	case elevio.MD_Up:
		return e.requests[floor][elevio.BT_HallUp] ||
			e.requests[floor][elevio.BT_Cab] ||
			!requestsAbove(e)

	case elevio.MD_Down:
		return e.requests[floor][elevio.BT_HallDown] ||
			e.requests[floor][elevio.BT_Cab] ||
			!requestsBelow(e)

	case elevio.MD_Stop:
		fallthrough
	default:
		return true
	}
}

func clearRequestAtFloor(elevator *Elevator, orderCompleteCh chan<- elevio.ButtonEvent) {
	for button := 0; button < config.N_BUTTONS; button++ {
		if button != elevio.BT_Cab && elevator.requests[elevator.floor][button] {
			orderCompleteCh <- elevio.ButtonEvent{Floor: elevator.floor, Button: elevio.ButtonType(button)}
		}
		elevator.requests[elevator.floor][button] = false
	}

}

func clearRequestAtFloorSimulation(elevator *Elevator) {
	for button := 0; button < config.N_BUTTONS; button++ {
		elevator.requests[elevator.floor][button] = false
	}
}

func clearAllRequest(elevator *Elevator) {
	for floor := 0; floor < config.N_FLOORS; floor++ {
		for button := 0; button < config.N_BUTTONS; button++ {
			elevator.requests[floor][button] = false
		}
	}
}

func clearAllHallRequests(elevator *Elevator) {
	for floor := 0; floor < config.N_FLOORS; floor++ {
		elevator.requests[floor][elevio.BT_HallDown] = false
		elevator.requests[floor][elevio.BT_HallUp] = false
	}
}
