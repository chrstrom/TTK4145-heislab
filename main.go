package main

import (
	"./elevio"
	"./fsm"
	"math/rand"
	"time"

	"github.com/TTK4145-Students-2021/project-gruppe80/elevio"
	"github.com/TTK4145-Students-2021/project-gruppe80/localOrderDelegation"
	"github.com/TTK4145-Students-2021/project-gruppe80/mock"
	"github.com/TTK4145-Students-2021/project-gruppe80/network"
	"github.com/TTK4145-Students-2021/project-gruppe80/orderDelegation"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	//numFloors := 4

	//elevio.Init("localhost:15657", numFloors)

	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go fsm.RunElevatorFSM(drv_buttons, drv_floors, drv_obstr, drv_stop)

	requestToNetwork := make(chan network.NewRequest)
	delegateOrderToNetwork := make(chan network.Delegation)
	requestReplyToNetwork := make(chan network.RequestReply)
	delegationConfirmToNetwork := make(chan network.DelegationConfirm)
	orderCompleteToNetwork := make(chan network.OrderComplete)

	requestFromNetwork := make(chan network.NewRequest)
	delegateFromNetwork := make(chan network.Delegation)
	requestReplyFromNetwork := make(chan network.RequestReply)
	delegationComfirmFromNetwork := make(chan network.DelegationConfirm)
	orderCompleteFromNetwork := make(chan network.OrderComplete)

	var node network.Node

	id := node.Init(requestToNetwork, delegateOrderToNetwork, requestReplyToNetwork, delegationConfirmToNetwork, orderCompleteToNetwork,
		requestFromNetwork, delegateFromNetwork, requestReplyFromNetwork, delegationComfirmFromNetwork, orderCompleteFromNetwork)
	go node.NetworkNode()

	cabOrderChannel := make(chan int)
	hallOrderChannel := make(chan localOrderDelegation.LocalOrder)

	var localOrderDelegator localOrderDelegation.LocalOrderDelegator
	localOrderDelegator.Init(drv_buttons, cabOrderChannel, hallOrderChannel)
	go localOrderDelegator.LocalOrderDelegation()

	var orderDelegator orderDelegation.OrderDelegator
	orderDelegator.Init(id, hallOrderChannel, requestToNetwork, delegateOrderToNetwork, requestReplyFromNetwork, delegationComfirmFromNetwork)
	go orderDelegator.OrderDelegation()

	/** 	mock functions for testing 		**/
	go mock.ReplyToRequests(requestFromNetwork, requestReplyToNetwork)
	go mock.ReplyToDelegations(delegateFromNetwork, delegationConfirmToNetwork)
	go mock.SendButtonPresses(drv_buttons, time.Second*10)

	//o := network.NewRequest{OrderID: 1, Floor: 1, Dir: 0}
	for {
		time.Sleep(time.Second * 5)
		//requestToNetwork <- o
		//o.OrderID++
	}
	/*for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}*/
}
