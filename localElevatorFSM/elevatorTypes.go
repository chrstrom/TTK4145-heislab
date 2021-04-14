package fsm

import (
	io "../elevio"
)

const N_FLOORS = 4
const N_BUTTONS = 3

//Not final values, only for efficient testing
const DOOR_OPEN_DURATION = 2
const TRAVEL_TIME = 1

type ElevatorState int

const (
	DoorOpen ElevatorState = iota
	Moving
	Idle
)

type Elevator struct {
	floor        int
	direction    io.MotorDirection
	requests     [N_FLOORS][N_BUTTONS]bool
	state        ElevatorState
	timerChannel chan int
	timerResets  int
	obstruction  bool
}
