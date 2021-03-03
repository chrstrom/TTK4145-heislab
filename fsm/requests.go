package fsm

func requestsAbove() bool {
	for floor:=elevator.floor+1; i < N_FLOORS; floor++ {
		for button:=0; button < N_BUTTONS; button++ {
			if elevator.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func requestsBelow() bool {
	for floor:=0; i < elevator.floor; floor++ {
		for button:=0; button < N_BUTTONS; button++ {
			if elevator.requests[floor][button] {
				return true
			}
		}
	}
	return false
}

func chooseDirection() MotorDirection {
	switch elevator.direction {

	// Note which functions are called first for each case!
	// For MD_Stop it doesn't really matter which one goes first
	case MD_UP:
		if requestsAbove() {
			return MD_UP
		} else if requestsBelow() {
			return MD_Down
		} else {
			return MD_Stop
		}

	case MD_Stop:
		fallthrough

	case MD_Down:
		if requestsBelow() {
			return MD_Down
		} else if requestsAbove() {
			return MD_Up
		} else {
			return MD_Stop
		}
	}
}

func shouldStop() bool {

	var floor = elevator.floor

	switch elevator.direction {

	case MD_Up:
		return elevator.requests[floor][BT_HallUp] || 
			   elevator.requests[floor][BT_Cab] ||
			   !requestsAbove()

	case MD_Down:
		return elevator.requests[floor][BT_HallDown] || 
			   elevator.requests[floor][BT_Cab] ||
			   !requestsBelow()

	case MD_Stop:
		fallthrough
	default:
		return true
	}
}

func clearRequestAtFloor() {
	for button:=0; button < N_BUTTONS; button++ {
		elevator.requests[elevator.floor][button] = 0
}