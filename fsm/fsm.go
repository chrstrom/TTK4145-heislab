package fsm

// Orders are received from the elevator server, so it makes
// sense for the switch/cases to be replaced with a for/select
// 
// The functions would then also take in a channel, defined in main.go
// Internal state will change as normal by setting something equals
// to something else, but the outputs should also go on channels
// This ensures that the FSM's responsibility is limited.
// This also means that the "elevator output device" will not exist
// Outputs instead reach the elevator through the elevator server

// In this implementation the output channels are replaced by hardware actions
// For the actual 



func makeUninitializedElevator() *Elevator {
	elevator := new(Elevator)
	elevator.floor = -1
	elevator.direction = Stop
	// 2D array of requests is 0 by default
	elevator.state = Idle

	return elevator
}




func onRequestButtonPress(button_msg ButtonEvent) {

	button_floor = button_msg.Floor
	button_type = button_msg.Button

	switch elevator.state {

	case DoorOpen:
		if(elevator.floor = button_floor) {
			timer_start(DOOR_OPEN_DURATION)
		}
		else {
			ele
			
			vator.requests[button_floor][button_type] = 1
		}

	}
}


// This function is the function from the fsm package that will run
// as a goroutine. Because of this, it should take inputs based on
// channels, and the for-select will take care of the 
//
//
//
func RunElevatorFSM(atFloor       FloorArrival <-chan,
		    buttonPress   ButtonPress  <-chan)
{

	// Initialize
	elevator = makeUninitializedElevator()

	// Loops indefinitely. RunElevatorFSM is a goroutine.
	for {

		select {

		case newFloor := <-atFloor:
			onFloorArrival(newFloor)
			
		case newButtonPress := <-buttonPress:
			onRequestButtonPress(newButtonPress)
		}

	}


}
