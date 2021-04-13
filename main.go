package main

import (
	"math/rand"
	"time"

	io "./elevio"
	"./hallOrderManager"
	fsm "./localElevatorFSM"
	"./localOrderDelegation"
	"./mock"
	"./network"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	testOrderManager()

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
	cabOrderChannel := make(chan int)
	hallOrderChannel := make(chan localOrderDelegation.LocalOrder)
	fsmChannels := fsm.CreateFSMChannelStruct()
	//delegateHallOrderChannel := make(chan io.ButtonEvent)
	//costChannel := make(chan int)
	//requestCostChannel := make(chan io.ButtonEvent)
	//orderCompleteChannel := make(chan io.ButtonEvent)

	// IO //
	go io.PollButtons(drv_buttons)
	go io.PollFloorSensor(drv_floors)
	go io.PollObstructionSwitch(drv_obstr)
	go io.PollStopButton(drv_stop)

	// Network //
	id := network.GetNodeID()
	go network.NetworkNode(id, networkChannels)

	// Elevator //
	go localOrderDelegation.OrderDelegator(drv_buttons, cabOrderChannel, hallOrderChannel)
	go hallOrderManager.OrderManager(id, hallOrderChannel, fsmChannels, networkChannels)
	go fsm.RunElevatorFSM(cabOrderChannel, fsmChannels, drv_floors, drv_obstr, drv_stop, timer_ch)

	for {
	}

}

func testOrderManager() {
	// Order Manager Channels //
	networkChannels := network.CreateNetworkChannelStruct()
	//cabOrderChannel := make(chan int)
	hallOrderChannel := make(chan localOrderDelegation.LocalOrder)
	fsmChannels := fsm.CreateFSMChannelStruct()
	// delegateHallOrderChannel := make(chan io.ButtonEvent)
	// costChannel := make(chan int)
	// requestCostChannel := make(chan io.ButtonEvent)
	// orderCompleteChannel := make(chan io.ButtonEvent)

	// Network //
	id := network.GetNodeID()
	go network.NetworkNode(id, networkChannels)

	// Elevator //
	go hallOrderManager.OrderManager(id, hallOrderChannel, fsmChannels, networkChannels)

	// /** 	mock functions for testing 		**/
	go mock.ReplyToRequests(networkChannels.RequestFromNetwork, networkChannels.ReplyToRequestToNetwork)
	go mock.Receive(fsmChannels.DelegateHallOrder)
	go mock.ElevatorCost(fsmChannels.RequestCost, fsmChannels.Cost)
	//go mock.ReplyToDelegations(networkChannels.DelegateFromNetwork, networkChannels.DelegationConfirmToNetwork)

	o := localOrderDelegation.LocalOrder{Floor: 2, Dir: 1}
	for {
		time.Sleep(time.Second * 5)
		hallOrderChannel <- o
	}
}
