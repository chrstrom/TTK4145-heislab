package main

import (
	"flag"
	"math/rand"
	"time"

	"./config"
	"./elevio"
	"./hallOrderManager"
	"./localElevatorFSM"
	"./localOrderDelegation"
	"./network"
)

func main() {
	elevatorPort := flag.String("port", "15657", "port number of the elevator server")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	elevio.Init("localhost:"+*elevatorPort, config.N_FLOORS)

	// IO Channels //
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	// Order Manager Channels //
	networkChannels := network.CreateNetworkChannelStruct()
	cabOrderChannel := make(chan int)
	hallOrderChannel := make(chan localOrderDelegation.LocalOrder)
	fsmChannels := localElevatorFSM.CreateFSMChannelStruct()

	// Hardware //
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Elevator // 
	id := network.GetNodeID()
	go network.NetworkNode(id, fsmChannels, networkChannels)
	go hallOrderManager.OrderManager(id, hallOrderChannel, fsmChannels, networkChannels)

	go localOrderDelegation.OrderDelegator(drv_buttons, cabOrderChannel, hallOrderChannel)
	go localElevatorFSM.RunElevatorFSM(cabOrderChannel, fsmChannels, drv_floors, drv_obstr, drv_stop, timer_ch)

	for {
		time.Sleep(time.Second * 10)
	}

}
