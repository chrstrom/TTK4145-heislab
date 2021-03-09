package fsm

import (
	io "../elevio"
)

const N_FLOORS = 4
const N_BUTTONS = 3

const DOOR_OPEN_DURATION = 5

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
	timerReset   int
	obstruction  bool
}
