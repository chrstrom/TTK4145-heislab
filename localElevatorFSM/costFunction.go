package fsm

import (
	io "../elevio"
)

func timeToIdle(e Elevator, floor int, button int) int {
	var duration int = 0
	e.requests[floor][button] = true

	switch e.state {

	case Idle:
		e.direction = chooseDirection(e)
		if e.direction == io.MD_Stop {
			return duration
		}
		break
	case Moving:
		duration = TRAVEL_TIME / 2
		e.floor += int(e.direction)
		break
	case DoorOpen:
		duration -= DOOR_OPEN_DURATION / 2
	}

	for {
		if shouldStop(e) {
			e = clearRequestAtFloorSimulation(e)
			if e.floor == floor {
				return duration
			}
			duration += DOOR_OPEN_DURATION
			e.direction = chooseDirection(e)

		}
		e.floor += int(e.direction)
		duration += TRAVEL_TIME
	}
}
