package fsm

type Button int
const (
	Cab Button = iota
	Floor
)

type Direction int
const (
	Up Direction = iota
	Down 
	Stop
)

type ElevatorState int
const (
	DoorOpen ElevatorState = iota
	Moving
	Idle
)

type Elevator struct {
	floor int
	direction Direction
	requests [N_FLOORS][N_BUTTONS]int
	state ElevatorState
}


