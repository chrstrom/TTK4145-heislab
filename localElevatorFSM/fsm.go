package fsm

import (
	"fmt"
	"time"

	"../cabOrderStorage"
	io "../elevio"
	types "../orderTypes"
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
var elevatorSimulator = makeUninitializedElevator()
var cost = 0

func CreateFSMChannelStruct() types.FSMChannels {
	var fsmChannels types.FSMChannels

	fsmChannels.DelegateHallOrder = make(chan io.ButtonEvent)
	fsmChannels.Cost = make(chan int)
	fsmChannels.RequestCost = make(chan io.ButtonEvent)
	fsmChannels.OrderComplete = make(chan io.ButtonEvent)

	return fsmChannels
}

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
			elevator.direction = chooseDirection(elevator)
			io.SetMotorDirection(elevator.direction)
			elevator.state = Moving
		}
	}

	setCabLights()
}

func onFloorArrival(floor int, orderCompleteCh chan<- io.ButtonEvent) {
	elevator.floor = floor

	// Set floor light
	io.SetFloorIndicator(floor)

	switch elevator.state {

	case Moving:

		if shouldStop(elevator) {
			io.SetMotorDirection(io.MD_Stop)
			elevator = clearRequestAtFloor(elevator, orderCompleteCh)

			doorOpenTimer()
			setCabLights()
			
			elevator.state = DoorOpen
		}

	}
}

func onDoorTimeout() {
	if elevator.state == DoorOpen && !elevator.obstruction {
		io.SetDoorOpenLamp(false)
		elevator.direction = chooseDirection(elevator)
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

func setCabLights() {
	cab_button := io.ButtonType(2)

	for f := 0; f < N_FLOORS; f++ {

		if elevator.requests[f][cab_button] {
			io.SetButtonLamp(cab_button, f, true)
		} else {
			io.SetButtonLamp(cab_button, f, false)
		}

	}
}

func doorOpenTimer() {
	const doorOpenTime = time.Second * DOOR_OPEN_DURATION
	io.SetDoorOpenLamp(true)
	timer.FsmSendWithDelay(doorOpenTime, elevator.timerChannel)
}

// This function is the function from the fsm package that will run
// as a goroutine. Because of this, it should take inputs based on
// channels, and the for-select will take care of the
func RunElevatorFSM(event_cabOrder <-chan int,
	fsmChannels types.FSMChannels,
	event_floorArrival <-chan int,
	event_obstruction <-chan bool,
	event_stopButton <-chan bool,
	event_timer <-chan int) {

	//Load cab ordes
	cabOrders := cabOrderStorage.LoadCabOrders()
	for f := 0; f < N_FLOORS; f++ {
		elevator.requests[f][2] = cabOrders[f]
	}
	fmt.Printf("CabOrders loaded!\n")
	fmt.Printf("CabOrders %+v\n", cabOrders)

	if elevator.floor == -1 {
		onInitBetweenFloors()
	}

	// Loops indefinitely. RunElevatorFSM *should be* a goroutine.
	for {
		// Dooropen=0, Moving=1, Idle=2
		fmt.Printf("State:%+v\n", elevator.state)

		//Store cab orders
		cabOrderStorage.StoreCabOrders(elevator.requests)

		select {

		case cabOrder := <-event_cabOrder:
			fmt.Printf("%+v\n", cabOrder)
			onRequestButtonPress(io.ButtonEvent{Floor: cabOrder, Button: io.BT_Cab})

		case costRequested := <-fsmChannels.RequestCost:
			elevatorSimulator = elevator
			cost = timeToIdle(elevatorSimulator, costRequested.Floor, int(costRequested.Button))
			fsmChannels.Cost <- cost

		case delegatedHallOrder := <-fsmChannels.DelegateHallOrder:
			fmt.Printf("Hallorder recieved!\n")
			onRequestButtonPress(delegatedHallOrder)

		case newFloor := <-event_floorArrival:
			fmt.Printf("%+v\n", newFloor)
			onFloorArrival(newFloor, fsmChannels.OrderComplete)

		case obstruction := <-event_obstruction:
			if obstruction {
				elevator.obstruction = true
				fmt.Printf("Obstruction triggered!\n")
			} else {
				elevator.obstruction = false
				fmt.Printf("No obstruction\n")
			}
			onObstruction(elevator.obstruction)

		/*case emergencyStop := <-event_stopButton:
		if emergencyStop {
			fmt.Printf("Emergency stop button triggered!")
		}
		//onEmergencyStop(emergencyStop)*/

		case timer := <-elevator.timerChannel:
			elevator.timerResets += timer
			if elevator.timerResets == 0 {
				onDoorTimeout()
			}

		}

	}
}
