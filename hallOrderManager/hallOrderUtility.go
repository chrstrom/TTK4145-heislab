package hallOrderManager

import (
	"log"
	"os"

	"../config"
	"../elevio"
	"../localOrderDelegation"
	msg "../orderTypes"
)

func getIDOfLowestCost(costs map[string]int, defaultID string) string {
	lowest := 100000000
	lowestID := ""

	for id, c := range costs {
		if c <= lowest {
			lowest = c
			lowestID = id
		}
	}

	if lowestID == "" {
		lowestID = defaultID
	}
	
	return lowestID
}

func setHallLight(dir int, floor int, state bool) {
	elevio.SetButtonLamp(elevio.ButtonType(dir), floor, state)
}

func initializeManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	fsmChannels msg.FSMChannels,
	channels msg.NetworkChannels) HallOrderManager {

	var manager HallOrderManager

	manager.id = id

	manager.orders = make(OrderMap)
	manager.orderIDCounter = 1

	manager.localRequestChannel = localRequestCh

	manager.requestToNetwork = channels.RequestToNetwork
	manager.delegateToNetwork = channels.DelegateOrderToNetwork
	manager.orderSyncToNetwork = channels.SyncOrderToNetwork
	manager.delegationConfirmToNetwork = channels.DelegationConfirmToNetwork

	manager.replyToRequestFromNetwork = channels.ReplyToRequestFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork
	manager.orderSyncFromNetwork = channels.SyncOrderFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork
	manager.delegationFromNetwork = channels.DelegateFromNetwork
	manager.peerUpdateChannel = channels.PeerUpdate

	manager.delegateToLocalElevator = fsmChannels.DelegateHallOrder
	manager.elevatorCost = fsmChannels.ReplyToHallOrderManager
	manager.requestElevatorCost = fsmChannels.RequestCost
	manager.orderComplete = fsmChannels.OrderComplete

	manager.orderReplyTimeoutChannel = make(chan int)
	manager.orderDelegationTimeoutChannel = make(chan int)
	manager.orderCompleteTimeoutChannel = make(chan msg.HallOrder)

	filepath := "log/" + manager.id + "-hallOrderManager.log"
	file, _ := os.Create(filepath)
	manager.logger = log.New(file, "", log.Ltime|log.Lmicroseconds)

	// Turn off all hall lights on init
	for f := 0; f < config.N_FLOORS; f++ {
		for b := elevio.ButtonType(0); b < 2; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}

	return manager
}
