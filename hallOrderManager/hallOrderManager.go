package hallOrderManager

import (
	"fmt"
	"log"

	"os"

	"../elevio"
	"../localOrderDelegation"
	"../network/peers"
	msg "../orderTypes"
	"../timer"
	"../utility"
)

func OrderManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	fsmChannels msg.FSMChannels,
	channels msg.NetworkChannels) {

	manager := initializeManager(id, localRequestCh, fsmChannels, channels)

	for {
		select {
		case request := <-manager.localRequestChannel:
			handleLocalRequest(request, &manager)
		case orderComplete := <-manager.orderComplete:
			//Get orders at the same floor from ordermap
			//Send confirmation to network
			fmt.Printf("Order %v completed\n", orderComplete)

		case reply := <-manager.replyToRequestFromNetwork:
			handleReplyFromNetwork(reply, &manager)

		case confirm := <-manager.orderDelegationConfirmFromNetwork:
			handleConfirmationFromNetwork(confirm, &manager)

		case delegation := <-manager.delegationFromNetwork:
			acceptDelegatedHallOrder(delegation, &manager)

		case order := <-manager.orderSyncFromNetwork:
			synchronizeOrderFromNetwork(order, &manager)

		case orderID := <-manager.orderReplyTimeoutChannel:
			delegateHallOrder(orderID, &manager)

		case orderID := <-manager.orderDelegationTimeoutChannel:
			selfServeHallOrder(orderID, &manager)

		case order := <-manager.orderCompleteTimeoutChannel:
			handleOrderCompleteTimeout(order, &manager)

		case peerUpdate := <-manager.peerUpdateChannel:
			handlePeerUpdate(peerUpdate, &manager)

			//case <-time.After(time.Second * 5):
			//	manager.logger.Printf("Quiet for 5 seconds")
		}
	}
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
	manager.elevatorCost = fsmChannels.Cost
	manager.requestElevatorCost = fsmChannels.RequestCost
	manager.orderComplete = fsmChannels.OrderComplete

	manager.orderReplyTimeoutChannel = make(chan int)
	manager.orderDelegationTimeoutChannel = make(chan int)
	manager.orderCompleteTimeoutChannel = make(chan msg.HallOrder)

	filepath := "log/" + manager.id + "-hallOrderManager.log"
	file, _ := os.Create(filepath)
	manager.logger = log.New(file, "", log.Ltime|log.Lmicroseconds)

	return manager
}

func handleLocalRequest(request localOrderDelegation.LocalOrder, manager *HallOrderManager) {
	//Check if order already exits? Or is this better to do in localOrdermanager? Or allow duplicates

	// This order will get synced with every elevator on the network
	order := msg.HallOrder{
		OwnerID: manager.id,
		ID:      manager.orderIDCounter,
		State:   msg.Received,
		Floor:   request.Floor,
		Dir:     request.Dir}

	manager.orderIDCounter++
	order.Costs = make(map[string]int)

	//get local elevator cost in some way
	manager.requestElevatorCost <- elevio.ButtonEvent{Floor: order.Floor, Button: elevio.ButtonType(order.Dir)}
	order.Costs[manager.id] = <-manager.elevatorCost * 0
	fmt.Printf("Cost:%v\n", order.Costs[manager.id])
	//order.Costs[manager.id] = rand.Intn(1000)

	manager.orders.update(order)

	//fmt.Printf("%v - local request received \n", order.ID)
	timer.SendWithDelay(orderReplyTime, manager.orderReplyTimeoutChannel, order.ID)
	timer.SendWithDelayHallOrder(orderCompletionTimeout, manager.orderCompleteTimeoutChannel, order)

	orderToNet := msg.OrderStamped{
		OrderID: order.ID,
		Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}

	manager.logger.Printf("New order ID%v: %#v", order.ID, order)
	manager.requestToNetwork <- orderToNet
}

func handleReplyFromNetwork(reply msg.OrderStamped, manager *HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, reply.OrderID)
	if valid && order.State == msg.Received {
		order.Costs[reply.ID] = reply.Order.Cost
		manager.orders.update(order)
		manager.logger.Printf("New reply to order ID%v: %#v", order.ID, order)
	}
}

func handleConfirmationFromNetwork(confirm msg.OrderStamped, manager *HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, confirm.OrderID)
	if valid && order.State == msg.Delegate {
		order.State = msg.Serving
		//fmt.Printf("%v - delegation confirmed \n", confirm.OrderID)

		manager.orders.update(order)

		manager.logger.Printf("Confirmed ID%v: %#v", order.ID, order)
		//Let the elevators on the network know that this local elevator has taken an order
		orderStateBroadcast(order, manager)
	}
}

func acceptDelegatedHallOrder(delegation msg.OrderStamped, manager *HallOrderManager) {
	order := msg.HallOrder{OwnerID: delegation.ID,
		ID:            delegation.OrderID,
		DelegatedToID: manager.id,
		State:         msg.Serving,
		Floor:         delegation.Order.Floor,
		Dir:           delegation.Order.Dir}

	manager.orders.update(order)
	manager.logger.Printf("Received order from net: %#v", order)
	manager.delegateToLocalElevator <- elevio.ButtonEvent{Floor: delegation.Order.Floor, Button: elevio.ButtonType(delegation.Order.Dir)}

	reply := msg.OrderStamped{
		ID:      order.OwnerID,
		OrderID: order.ID,
		Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}
	manager.delegationConfirmToNetwork <- reply
}

func synchronizeOrderFromNetwork(order msg.HallOrder, manager *HallOrderManager) {
	// Receive an order from the network and add it to the list of hall orders
	if order.OwnerID != manager.id {
		manager.orders.update(order)
		manager.logger.Printf("Sync from net: %#v", order)
	}
}

func delegateHallOrder(orderID int, manager *HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, orderID)
	if valid && order.State == msg.Received {
		id := getIDOfLowestCost(order.Costs)
		if id == "" {
			id = manager.id
		}
		order.DelegatedToID = id

		if id == manager.id {
			//send order to local elevator
			manager.delegateToLocalElevator <- elevio.ButtonEvent{Floor: order.Floor, Button: elevio.ButtonType(order.Dir)}

			manager.logger.Printf("Delegate order ID%v to local elevator (%v replies): %#v", order.ID, len(order.Costs), order)
			order.State = msg.Serving
			orderStateBroadcast(order, manager)
		} else {
			manager.logger.Printf("Delegate order ID%v to net (%v replies): %#v", order.ID, len(order.Costs), order)
			//fmt.Printf("%v - delegate to %v  (%v replies) \n", orderID, id, len(order.Costs))
			timer.SendWithDelay(orderDelegationTime, manager.orderDelegationTimeoutChannel, orderID)

			order.State = msg.Delegate

			message := msg.OrderStamped{
				ID:      order.DelegatedToID,
				OrderID: orderID,
				Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}

			manager.delegateToNetwork <- message
		}
		manager.orders.update(order)
	}
}

func selfServeHallOrder(orderID int, manager *HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, orderID)
	if valid && order.State == msg.Delegate {
		//Send order to local elevator
		manager.delegateToLocalElevator <- elevio.ButtonEvent{Floor: order.Floor, Button: elevio.ButtonType(order.Dir)}
		order.DelegatedToID = manager.id
		order.State = msg.Serving

		//fmt.Printf("------------- %v - delegation timedout! Sending to local elevator \n", orderID)
		manager.logger.Printf("Timeout delegation ID%v, sending to local elevator: %v", order.ID, order)

		manager.orders.update(order)
		orderStateBroadcast(order, manager)
	}
}

func orderStateBroadcast(order msg.HallOrder, manager *HallOrderManager) {
	manager.orderSyncToNetwork <- order
	manager.logger.Printf("Sync order ID%v to net:%#v", order.ID, order)
}

func handleOrderCompleteTimeout(order msg.HallOrder, manager *HallOrderManager) {
	manager.logger.Printf("Order timeout ID%v: %#v", order.ID, order)
	redelegateOrder(order, manager)
}

func redelegateOrder(o msg.HallOrder, manager *HallOrderManager) {
	order, ok := manager.orders.getOrder(o.OwnerID, o.ID)
	if ok && order.State != msg.Completed {
		manager.logger.Printf("Redelegate order ID%v: %#v", order.ID, order)
		if order.OwnerID != manager.id {
			order.OwnerID = manager.id
			order.ID = manager.orderIDCounter
			manager.orderIDCounter++
		}
		order.Costs = make(map[string]int)
		order.DelegatedToID = ""
		order.State = msg.Received

		manager.requestElevatorCost <- elevio.ButtonEvent{Floor: order.Floor, Button: elevio.ButtonType(order.Dir)}
		order.Costs[manager.id] = <-manager.elevatorCost

		manager.orders.update(order)

		timer.SendWithDelay(orderReplyTime, manager.orderReplyTimeoutChannel, order.ID)
		timer.SendWithDelayHallOrder(orderCompletionTimeout, manager.orderCompleteTimeoutChannel, order)

		orderToNet := msg.OrderStamped{
			OrderID: order.ID,
			Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}

		manager.requestToNetwork <- orderToNet
	}
}

func handlePeerUpdate(peerUpdate peers.PeerUpdate, manager *HallOrderManager) {
	for _, nodeid := range peerUpdate.Lost {
		manager.logger.Printf("Node lost connection: %v", nodeid)
		orders := manager.orders.getOrderToID(nodeid)
		for _, o := range orders {
			if o.OwnerID == manager.id {
				redelegateOrder(o, manager)
				fmt.Println("redelegeate 1")
			} else if !utility.IsStringInSlice(o.OwnerID, peerUpdate.Peers) && !utility.IsStringInSlice(o.DelegatedToID, peerUpdate.Peers) {
				redelegateOrder(o, manager)
				fmt.Println("redelegeate 2")
			}
		}
	}

	if len(peerUpdate.New) > 0 {
		manager.logger.Printf("New node(s) connected")
		orders := manager.orders.getOrderFromID(manager.id)
		for _, o := range orders {
			if o.State == msg.Serving {
				orderStateBroadcast(o, manager)
			}
		}
	}
}
