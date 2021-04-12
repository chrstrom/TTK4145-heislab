package fsm

import io "../elevio"

func requestsAbove(e Elevator) bool {
	for floor := e.floor + 1; floor < N_FLOORS; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			if e.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for floor := 0; floor < e.floor; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			if e.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func chooseDirection(e Elevator) io.MotorDirection {
	switch e.direction {

	// Note which functions are called first for each case!
	// For MD_Stop it doesn't really matter which one goes first
	case io.MD_Up:
		if requestsAbove(e) {
			return io.MD_Up
		} else if requestsBelow(e) {
			return io.MD_Down
		} else {
			return io.MD_Stop
		}

	case io.MD_Stop:

		fallthrough

	case io.MD_Down:
		if requestsBelow(e) {
			return io.MD_Down
		} else if requestsAbove(e) {
			return io.MD_Up
		} else {
			return io.MD_Stop
		}

	}
	return io.MD_Stop
}

func shouldStop(e Elevator) bool {

	var floor = e.floor

	switch e.direction {

	case io.MD_Up:
		return e.requests[floor][io.BT_HallUp] ||
			e.requests[floor][io.BT_Cab] ||
			!requestsAbove(e)

	case io.MD_Down:
		return e.requests[floor][io.BT_HallDown] ||
			e.requests[floor][io.BT_Cab] ||
			!requestsBelow(e)

	case io.MD_Stop:
		fallthrough
	default:
		return true
	}
}

func clearRequestAtFloor(e Elevator, orderCompleteCh chan<- io.ButtonEvent) Elevator {
	for button := 0; button < N_BUTTONS; button++ {
		if button != io.BT_Cab && e.requests[e.floor][button] {
			orderCompleteCh <- io.ButtonEvent{Floor: e.floor, Button: io.ButtonType(button)}
		}
		e.requests[e.floor][button] = false
	}
	return e
}

func clearRequestAtFloorSimulation(e Elevator) Elevator {
	for button := 0; button < N_BUTTONS; button++ {
		e.requests[e.floor][button] = false
	}
	return e
}

func clearAllRequest(e Elevator) Elevator {
	for floor := 0; floor < N_FLOORS; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			e.requests[floor][button] = false
		}
	}
	return e
}
