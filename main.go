package main

import (
	"math/rand"
	"time"

	io "./elevio"
	"./fsm"
	"./hallOrderManager"
	"./localOrderDelegation"
	"./network"
)

func main() {
	hallOrderManager.TestOrderMap()
	return

	rand.Seed(time.Now().UTC().UnixNano())

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
	localRequestChannel := make(<-chan localOrderDelegation.LocalOrder)
	//cabOrderChannel  := make(chan int)
	//hallOrderChannel := make(chan localOrderDelegation.LocalOrder)

	// IO //
	go io.PollButtons(drv_buttons)
	go io.PollFloorSensor(drv_floors)
	go io.PollObstructionSwitch(drv_obstr)
	go io.PollStopButton(drv_stop)

	// Network //
	id := network.GetNodeID()
	go network.NetworkNode(id, networkChannels)

	// Elevator //
	go hallOrderManager.OrderManager(id, localRequestChannel, networkChannels)
	go fsm.RunElevatorFSM(drv_buttons, drv_floors, drv_obstr, drv_stop, timer_ch)

	for {
	}

	// /** 	mock functions for testing 		**/
	// go mock.ReplyToRequests(requestFromNetwork, requestReplyToNetwork)
	// go mock.ReplyToDelegations(delegateFromNetwork, delegationConfirmToNetwork)
	// go mock.SendButtonPresses(drv_buttons, time.Second*10)

	//o := network.NewRequest{OrderID: 1, Floor: 1, Dir: 0}
	// for {
	// 	time.Sleep(time.Second * 5)
	//requestToNetwork <- o
	//o.OrderID++
	//}
}
