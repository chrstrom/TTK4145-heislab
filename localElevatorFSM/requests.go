package fsm

import io "../elevio"

func requestsAbove() bool {
	for floor := elevator.floor + 1; floor < N_FLOORS; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			if elevator.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func requestsBelow() bool {
	for floor := 0; floor < elevator.floor; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			if elevator.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func chooseDirection() io.MotorDirection {
	switch elevator.direction {

	// Note which functions are called first for each case!
	// For MD_Stop it doesn't really matter which one goes first
	case io.MD_Up:
		if requestsAbove() {
			return io.MD_Up
		} else if requestsBelow() {
			return io.MD_Down
		} else {
			return io.MD_Stop
		}

	case io.MD_Stop:

		fallthrough

	case io.MD_Down:
		if requestsBelow() {
			return io.MD_Down
		} else if requestsAbove() {
			return io.MD_Up
		} else {
			return io.MD_Stop
		}

	}
	return io.MD_Stop
}

func shouldStop() bool {

	var floor = elevator.floor

	switch elevator.direction {

	case io.MD_Up:
		return elevator.requests[floor][io.BT_HallUp] ||
			elevator.requests[floor][io.BT_Cab] ||
			!requestsAbove()

	case io.MD_Down:
		return elevator.requests[floor][io.BT_HallDown] ||
			elevator.requests[floor][io.BT_Cab] ||
			!requestsBelow()

	case io.MD_Stop:
		fallthrough
	default:
		return true
	}
}

func clearRequestAtFloor() {
	for button := 0; button < N_BUTTONS; button++ {
		elevator.requests[elevator.floor][button] = false
	}
}

func clearAllRequest() {
	for floor := 0; floor < N_FLOORS; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			elevator.requests[floor][button] = false
		}
	}
}
