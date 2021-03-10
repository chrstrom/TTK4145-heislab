package fsm

import (
	"fmt"
	"time"

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
	elevator.timerChannel = make(chan int)
	elevator.timerResets = 0
	elevator.obstruction = false

	return *elevator
}

func onInitBetweenFloors() {
	elevator.direction = io.MD_Down
	io.SetMotorDirection(elevator.direction)
	elevator.state = Moving
}

func onRequestButtonPress(button_msg io.ButtonEvent) {

	var button_floor = button_msg.Floor
	var button_type = button_msg.Button

	switch elevator.state {

	case DoorOpen:
		if elevator.floor == button_floor {
			doorOpenTimer()

		} else {
			elevator.requests[button_floor][button_type] = true
		}

	case Moving:
		elevator.requests[button_floor][button_type] = true

	case Idle:
		if elevator.floor == button_floor {
			io.SetDoorOpenLamp(true)
			doorOpenTimer()
			elevator.state = DoorOpen
		} else {
			elevator.requests[button_floor][button_type] = true
			elevator.direction = chooseDirection()
			io.SetMotorDirection(elevator.direction)
			elevator.state = Moving
		}
	}

	// TODO
	setAllLights()
}

func onFloorArrival(floor int) {
	elevator.floor = floor

	// Set floor light

	switch elevator.state {

	case Moving:
		// TODO
		if shouldStop() {
			io.SetMotorDirection(io.MD_Stop)
			clearRequestAtFloor()

			//TODO start timer with DOOR_OPEN_DURATION
			doorOpenTimer()

			// Set all order lights again
			setAllLights()
			elevator.state = DoorOpen
		}

	}
}

func onDoorTimeout() {
	if elevator.state == DoorOpen && !elevator.obstruction {
		//fmt.Printf("OnDoortimeout\n")
		io.SetDoorOpenLamp(false)
		elevator.direction = chooseDirection()
		io.SetMotorDirection(elevator.direction)

		if elevator.direction == io.MD_Stop {
			elevator.state = Idle
		} else {
			elevator.state = Moving
		}
	}
}

func onObstruction(obstruction bool) {
	if elevator.state == DoorOpen {
		io.SetDoorOpenLamp(true)
	}

	onDoorTimeout()
}

func onEmergencyStop(stop bool) {
	if stop {
		io.SetStopLamp(true)
		io.SetMotorDirection(io.MD_Stop)
		clearAllRequest()
		setAllLights()
	} else {
		io.SetStopLamp(false)
		if elevator.state == Moving {
			elevator.floor = -1
		}

	}
}

// This function is the function from the fsm package that will run
// as a goroutine. Because of this, it should take inputs based on
// channels, and the for-select will take care of the
func RunElevatorFSM(event_orderButton <-chan io.ButtonEvent,
	event_floorArrival <-chan int,
	event_obstruction <-chan bool,
	event_stopButton <-chan bool,
	event_timer <-chan int) {

	// Loops indefinitely. RunElevatorFSM *should be* a goroutine.

	for {

		if elevator.floor == -1 {
			onInitBetweenFloors()
		}

		for f := 0; f < N_FLOORS; f++ {
			for b := io.ButtonType(0); b < N_BUTTONS; b++ {
				if elevator.requests[f][b] && elevator.state == Idle {
					fmt.Printf("Looking through requests\n")
					onRequestButtonPress(io.ButtonEvent{Floor: f, Button: io.ButtonType(b)})
				}
			}
		}
		// Dooropen=0, Moving=1, Idle=2
		fmt.Printf("State:%+v\n", elevator.state)

		select {

		case newButtonPress := <-event_orderButton:
			fmt.Printf("%+v\n", newButtonPress)
			onRequestButtonPress(newButtonPress)

		case newFloor := <-event_floorArrival:
			fmt.Printf("%+v\n", newFloor)
			onFloorArrival(newFloor)

		case obstruction := <-event_obstruction:
			if obstruction {
				elevator.obstruction = true
				fmt.Printf("Obstruction triggered!\n")
			} else {
				elevator.obstruction = false
				fmt.Printf("No obstruction\n")
			}
			onObstruction(elevator.obstruction)

		case emergencyStop := <-event_stopButton:
			if emergencyStop {
				fmt.Printf("Emergency stop button triggered!")
			}
			onEmergencyStop(emergencyStop)
		case timer := <-elevator.timerChannel:
			elevator.timerResets += timer
			if elevator.timerResets == 0 {
				onDoorTimeout()
			}

		}

	}
}

func setAllLights() {
	for f := 0; f < N_FLOORS; f++ {
		for b := io.ButtonType(0); b < N_BUTTONS; b++ {
			if elevator.requests[f][b] {
				io.SetButtonLamp(b, f, true)
			} else {
				io.SetButtonLamp(b, f, false)
			}

		}
	}

}

func doorOpenTimer() {
	const doorOpenTime = time.Millisecond * 2000
	io.SetDoorOpenLamp(true)
	timer.FsmSendWithDelay(doorOpenTime, elevator.timerChannel)
}
