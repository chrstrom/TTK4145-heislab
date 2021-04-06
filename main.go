package main

import (
	"math/rand"
	"time"

	io "./elevio"
	"./fsm"
	"./hallOrderManager"
	"./localOrderDelegation"
	"./mock"
	"./network"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	//testOrderManager()

	numFloors := 4
	io.Init("localhost:15657", numFloors)

	// IO Channels //
	drv_buttons := make(chan io.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	timer_ch := make(chan int)

	// Order Manager Channels //
	networkChannels := network.CreateNetworkChannelStruct()
	//cabOrderChannel  := make(chan int)
	hallOrderChannel := make(chan localOrderDelegation.LocalOrder)

	// IO //
	go io.PollButtons(drv_buttons)
	go io.PollFloorSensor(drv_floors)
	go io.PollObstructionSwitch(drv_obstr)
	go io.PollStopButton(drv_stop)

	// Network //
	id := network.GetNodeID()
	go network.NetworkNode(id, networkChannels)

	// Elevator //
	go hallOrderManager.OrderManager(id, hallOrderChannel, networkChannels)
	go fsm.RunElevatorFSM(drv_buttons, drv_floors, drv_obstr, drv_stop, timer_ch)

	//for {
	//}

}

func testOrderManager() {
	// Order Manager Channels //
	networkChannels := network.CreateNetworkChannelStruct()
	//cabOrderChannel  := make(chan int)
	hallOrderChannel := make(chan localOrderDelegation.LocalOrder)

	// Network //
	id := network.GetNodeID()
	go network.NetworkNode(id, networkChannels)

	// Elevator //
	go hallOrderManager.OrderManager(id, hallOrderChannel, networkChannels)

	// /** 	mock functions for testing 		**/
	go mock.ReplyToRequests(networkChannels.RequestFromNetwork, networkChannels.RequestReplyToNetwork)
	//go mock.ReplyToDelegations(networkChannels.DelegateFromNetwork, networkChannels.DelegationConfirmToNetwork)

	o := localOrderDelegation.LocalOrder{Floor: 2, Dir: 1}
	for {
		time.Sleep(time.Second * 5)
		hallOrderChannel <- o
	}
}
