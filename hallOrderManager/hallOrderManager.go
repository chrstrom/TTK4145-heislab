package hallOrderManager

import (
	"fmt"
	"math/rand"

	"../localOrderDelegation"
	"../network"
	"../timer"
)

func OrderManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	channels network.NetworkChannels) {

	manager := initializeManager(id, localRequestCh, channels)

	for {
		select {
		case request := <-manager.localRequestChannel:
			handleLocalRequest(request, manager)

		case reply := <-manager.requestReplyFromNetwork:
			handleReplyFromNetwork(reply, manager)

		case confirm := <-manager.orderDelegationConfirmFromNetwork:
			handleConfirmationFromNetwork(confirm, manager)

		case delegation := <-manager.delegationFromNetwork:
			acceptDelegatedHallOrder(delegation, manager)

		case order := <-manager.orderSyncFromNetwork:
			synchronizeNetworkOrders(order, manager)

		case orderID := <-manager.orderReplyTimeoutChannel:
			delegateHallOrder(orderID, manager)

		case orderID := <-manager.orderDelegationTimeoutChannel:
			selfServeHallOrder(orderID, manager)
		}
	}
}

func initializeManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	channels network.NetworkChannels) HallOrderManager {

	var manager HallOrderManager

	manager.id = id

	manager.orders = make(OrderMap)
	manager.orderIDCounter = 1

	manager.localRequestChannel = localRequestCh

	manager.requestToNetwork = channels.RequestToNetwork
	manager.delegateToNetwork = channels.DelegateOrderToNetwork
	manager.orderSyncToNetwork = channels.SyncOrderToNetwork
	manager.delegationConfirmToNetwork = channels.DelegationConfirmToNetwork

	manager.requestReplyFromNetwork = channels.RequestReplyFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork
	manager.orderSyncFromNetwork = channels.SyncOrderFromNetwork
	manager.requestReplyFromNetwork = channels.RequestReplyFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork
	manager.delegationFromNetwork = channels.DelegateFromNetwork

	manager.orderReplyTimeoutChannel = make(chan int)
	manager.orderDelegationTimeoutChannel = make(chan int)

	return manager
}

func handleLocalRequest(request localOrderDelegation.LocalOrder, manager HallOrderManager) {
	//Check if order already exits? Or is this better to do in localOrdermanager? Or allow duplicates
	order := HallOrder{
		OwnerID: manager.id,
		ID:      manager.orderIDCounter,
		State:   Received,
		Floor:   request.Floor,
		Dir:     request.Dir}

	manager.orderIDCounter++
	order.costs = make(map[string]int)
	//get local elevator cost in some way
	order.costs[manager.id] = rand.Intn(1000)

	manager.orders.update(order)
	//fmt.Printf("%v - local request received \n", order.ID)

	timer.SendWithDelay(orderReplyTime, manager.orderReplyTimeoutChannel, order.ID)

	orderToNet := network.OrderStamped{
		OrderID: order.ID,
		Order:   network.Order{Floor: order.Floor, Dir: order.Dir}}

	manager.requestToNetwork <- orderToNet
}

func handleReplyFromNetwork(reply network.OrderStamped, manager HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, reply.OrderID)
	if valid && order.State == Received {
		order.costs[reply.ID] = reply.Order.Cost
	}
}

func handleConfirmationFromNetwork(confirm network.OrderStamped, manager HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, confirm.OrderID)
	if valid && order.State == Delegate {
		order.State = Serving
		fmt.Printf("%v - delegation confirmed \n", confirm.OrderID)

		manager.orders.update(order)

		//Let the elevators on the network know that this local elevator has taken an order
		orderToNet := network.OrderSync{OrderID: order.ID, Floor: order.Floor, Dir: order.Dir}
		manager.orderSyncToNetwork <- orderToNet
	}
}

func acceptDelegatedHallOrder(delegation network.OrderStamped, manager HallOrderManager) {
	order := HallOrder{OwnerID: delegation.ID,
		ID:            delegation.OrderID,
		DelegatedToID: manager.id,
		State:         Serving,
		Floor:         delegation.Order.Floor,
		Dir:           delegation.Order.Dir}

	manager.orders.update(order)
	reply := network.OrderStamped{
		ID: order.OwnerID, 
		OrderID: order.ID, 
		Order: network.Order{Floor: order.Floor, Dir: order.Dir}}
	manager.delegationConfirmToNetwork <- reply
}

func synchronizeNetworkOrders(order network.OrderSync, manager HallOrderManager) {
	if order.ID != manager.id {
		order := HallOrder{}
		manager.orders.update(order)
	}
}

func delegateHallOrder(orderID int, manager HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, orderID)
	if valid && order.State == Received {
		id := getIDOfLowestCost(order.costs)
		if id == "" {
			id = manager.id
		}
		order.DelegatedToID = id

		if id == manager.id {
			//send order to local elevator
			fmt.Printf("%v - delegate to local elevator (%v replies) \n", orderID, len(order.costs))

			order.State = Serving
		} else {
			fmt.Printf("%v - delegate to %v  (%v replies) \n", orderID, id, len(order.costs))
			timer.SendWithDelay(orderDelegationTime, manager.orderDelegationTimeoutChannel, orderID)

			order.State = Delegate

			message := network.OrderStamped{
				ID:      order.DelegatedToID,
				OrderID: orderID,
				Order:   network.Order{Floor: order.Floor, Dir: order.Dir}}

			manager.delegateToNetwork <- message
		}
		manager.orders.update(order)
	}
}

func selfServeHallOrder(orderID int, manager HallOrderManager) {
	order, valid := manager.orders.getOrder(manager.id, orderID)
	if valid && order.State == Delegate {
		//Send order to local elevator
		order.DelegatedToID = manager.id
		order.State = Serving

		fmt.Printf("%v - delegation timedout! Sending to local elevator \n", orderID)

		manager.orders.update(order)
	}
}
