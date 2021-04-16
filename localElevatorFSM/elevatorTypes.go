package fsm

import (
	io "../elevio"
	"../config"
)


type ElevatorState int

const (
	DoorOpen ElevatorState = iota
	Moving
	Idle
)

type Elevator struct {
	floor        int
	direction    io.MotorDirection
	requests     [config.N_FLOORS][config.N_BUTTONS]bool
	state        ElevatorState
	timerChannel chan int
	timerResets  int
	obstruction  bool
}
