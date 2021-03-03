package fsm

import (
	"fmt"
	io "../elevio"
	"../timer"
)

// Orders are received from the elevator server, so it makes
// sense for the switch/cases to be replaced with a for/select
//
// The functions would then also take in a channel, defined in main.go
// Internal state will change as normal by setting something equals
// to something else, but the outputs should also go on channels
// This ensures that the FSM's responsibility is limited.
// This also means that the "elevator output device" will not exist
// Outputs instead reach the elevator through the elevator server

// Making this variable available to the entire fsm package makes sense, since
// it makes us not have to pass the elevator around as an argument everywhere.
// The for-select block in RunElevatorFSM() will at all times only execute
// one case, which means that we will not run into concurrency problems here either
var elevator = makeUninitializedElevator()


func makeUninitializedElevator() Elevator {
	elevator := new(Elevator)
	elevator.floor = -1
	elevator.direction = io.MD_Stop
	// 2D array of requests is 0 by default
	elevator.state = Idle

	return *elevator
}




func onRequestButtonPress(button_msg io.ButtonEvent) {

	var button_floor = button_msg.Floor
	var button_type = button_msg.Button

	switch elevator.state {

	case DoorOpen:
		if elevator.floor == button_floor {
			//TODO start timer with DOOR_OPEN_DURATION 
		} else {
			elevator.requests[button_floor][button_type] = 1
		}
	
	case Moving:
		elevator.requests[button_floor][button_type] = 1

	case Idle:
		if elevator.floor == button_floor {
			io.SetDoorOpenLamp(true)
			//TODO start timer with DOOR_OPEN_DURATION 
			elevator.state = DoorOpen
		} else {
			elevator.requests[button_floor][button_type] = 1
			elevator.direction = chooseDirection()
			io.SetMotorDirection(elevator.direction)
			elevator.state = Moving
		}
	}

	// TODO
	setAllLights(elevator);
}

func onFloorArrival(floor int) {
	elevator.floor = floor

	// Set floor light

	switch elevator.state {

	case DoorOpen:
			// TODO
		if shouldStop() {
			io.SetMotorDirection(MD_Stop)
			io.SetDoorOpenLamp(true)
			clearRequestAtCurrentFloor(elevator);

			//TODO start timer with DOOR_OPEN_DURATION 
			
			// Set all order lights again
			elevator.state = DoorOpen
		}	
	}
}

func onDoorTimeout() {

	switch elevator.state {

	case DoorOpen:
		elevator.direction = chooseDirection()

		io.SetDoorOpenLamp(false)
		io.SetMotorDirection(elevator.direction)

		if elevator.direction == MD_Stop {
			elevator.state = Idle
		} else {
			elevator.state = Moving
		}
	}
}

// This function is the function from the fsm package that will run
// as a goroutine. Because of this, it should take inputs based on
// channels, and the for-select will take care of the
func RunElevatorFSM(event_orderButton 	<-chan io.ButtonEvent,
					event_floorArrival 	<-chan int,
					event_obstruction 	<-chan bool,
					event_stopButton 	<-chan bool) {

	// Loops indefinitely. RunElevatorFSM *should be* a goroutine.
	for {

		select {

		case newButtonPress := <-event_orderButton:
			onRequestButtonPress(newButtonPress)

		case newFloor := <-event_floorArrival:
			onFloorArrival(newFloor)

		case obstruction := <-event_obstruction:
			if obstruction {
				fmt.Printf("Obstruction triggered!")
			}

		case emergencyStop := <-event_stopButton:
			if emergencyStop {
				fmt.Printf("Emergency stop button triggered!")
			}
		}
	}
}
