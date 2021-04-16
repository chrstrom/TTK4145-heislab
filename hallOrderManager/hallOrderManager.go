package hallOrderManager

import (
	"../config"
	"../elevio"
	"../localOrderDelegation"
	msg "../orderTypes"
	"../timer"
	"../utility"
)

// This is the driver function for the hall order manager
// and contains a for-select, thus should be called as a goroutine.
func OrderManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	fsmChannels msg.FSMChannels,
	channels msg.NetworkChannels) {

	manager := initializeManager(id, localRequestCh, fsmChannels, channels)

	for {
		select {
		case request := <-manager.localRequestChannel:

			for _, order := range manager.orders.getOrdersToFloorWithDir(request.Floor, request.Dir) {
				if order.State != msg.Completed {
					break
				}
			}

			order := msg.HallOrder{
				OwnerID: manager.id,
				ID:      manager.orderIDCounter,
				State:   msg.Received,
				Floor:   request.Floor,
				Dir:     request.Dir}

			manager.orderIDCounter++
			order.Costs = make(map[string]int)

			manager.requestElevatorCost <- msg.RequestCost{
				Order: msg.OrderStamped{
					OrderID: order.ID,
					Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}},
				RequestFrom: msg.HallOrderManager}

			order.Costs[manager.id] = <-manager.elevatorCost
			manager.orders.update(order)

			timer.SendWithDelay(config.ORDER_REPLY_TIME, manager.orderReplyTimeoutChannel, order.ID)
			timer.SendWithDelayHallOrder(config.ORDER_COMPLETION_TIMEOUT, manager.orderCompleteTimeoutChannel, order)

			orderToNet := msg.OrderStamped{
				OrderID: order.ID,
				Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}

			manager.logger.Printf("New order ID%v: %#v", order.ID, order)
			manager.requestToNetwork <- orderToNet

		case buttonEvent := <-manager.orderComplete:

			dir := int(buttonEvent.Button)
			floor := buttonEvent.Floor

			setHallLight(dir, floor, false)

			for _, order := range manager.orders.getOrdersToFloorWithDir(floor, dir) {
				order.State = msg.Completed
				manager.orders.update(order)
				manager.logger.Printf("Order completed ID%v: %#v\n", order.ID, order)
				syncOrderWithOtherElevators(order, &manager)
			}

		case reply := <-manager.replyToRequestFromNetwork:

			order, orderIsValid := manager.orders.getOrder(manager.id, reply.OrderID)

			if orderIsValid && order.State == msg.Received {
				order.Costs[reply.ID] = reply.Order.Cost
				manager.orders.update(order)
				manager.logger.Printf("New reply to order ID%v: %#v", order.ID, order)
			}

		case confirm := <-manager.orderDelegationConfirmFromNetwork:

			order, orderIsValid := manager.orders.getOrder(manager.id, confirm.OrderID)

			if orderIsValid && order.State == msg.Delegate {
				order.State = msg.Serving
				manager.orders.update(order)
				manager.logger.Printf("Confirmed ID%v: %#v", order.ID, order)
				syncOrderWithOtherElevators(order, &manager)

				setHallLight(order.Dir, order.Floor, true)
			}

		case delegation := <-manager.delegationFromNetwork:

			incomingOrder := msg.HallOrder{
				OwnerID:       delegation.ID,
				ID:            delegation.OrderID,
				DelegatedToID: manager.id,
				State:         msg.Serving,
				Floor:         delegation.Order.Floor,
				Dir:           delegation.Order.Dir}

			manager.orders.update(incomingOrder)
			manager.logger.Printf("Received order from net: %#v", incomingOrder)

			orderForFSM := elevio.ButtonEvent{
				Floor:  delegation.Order.Floor,
				Button: elevio.ButtonType(delegation.Order.Dir)}

			manager.delegateToLocalElevator <- orderForFSM

			replyToNetwork := msg.OrderStamped{
				ID:      incomingOrder.OwnerID,
				OrderID: incomingOrder.ID,
				Order:   msg.Order{Floor: incomingOrder.Floor, Dir: incomingOrder.Dir}}
			manager.delegationConfirmToNetwork <- replyToNetwork

			setHallLight(incomingOrder.Dir, incomingOrder.Floor, true)

		case order := <-manager.orderSyncFromNetwork:

			orderSaved, orderExists := manager.orders.getOrder(order.OwnerID, order.ID)

			// Conditionally synchronize the order itself
			if !orderExists || (orderExists && order.State >= orderSaved.State) {

				if !orderExists {
					timer.SendWithDelayHallOrder(
						config.ORDER_COMPLETION_TIMEOUT,
						manager.orderCompleteTimeoutChannel,
						order)
				}
				manager.orders.update(order)
				manager.logger.Printf("Sync from net: %#v", order)

				// Make sure lights are updated in accordance with the synched order
				if order.State == msg.Serving {
					setHallLight(order.Dir, order.Floor, true)
				}

				if order.State == msg.Completed {
					setHallLight(order.Dir, order.Floor, false)
				}

			}

		case orderID := <-manager.orderReplyTimeoutChannel:
			// Once the window of time we listen to replies runs out,
			// we can start delegating the hall order

			order, orderIsValid := manager.orders.getOrder(manager.id, orderID)

			if orderIsValid && order.State == msg.Received {
				id := getIDOfLowestCost(order.Costs, manager.id)

				order.DelegatedToID = id

				if id == manager.id {
					selfServeHallOrder(order, &manager)
					order.State = msg.Serving

					manager.logger.Printf("Timeout delegation ID%v, sending to local elevator: %v", order.ID, order)

					syncOrderWithOtherElevators(order, &manager)
					setHallLight(order.Dir, order.Floor, true)

				} else { // Order will be given to another elevator on the network
					manager.logger.Printf("Delegate order ID%v to net (%v replies): %#v", order.ID, len(order.Costs), order)
					timer.SendWithDelay(config.ORDER_DELEGATION_TIME, manager.orderDelegationTimeoutChannel, orderID)

					order.State = msg.Delegate

					delegatedOrder := msg.OrderStamped{
						ID:      order.DelegatedToID,
						OrderID: orderID,
						Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}

					manager.delegateToNetwork <- delegatedOrder
				}
				manager.orders.update(order)
			}

		case orderID := <-manager.orderDelegationTimeoutChannel:

			order, orderIsValid := manager.orders.getOrder(manager.id, orderID)

			if orderIsValid && order.State == msg.Delegate {

				order.DelegatedToID = manager.id
				selfServeHallOrder(order, &manager)
				order.State = msg.Serving

				manager.logger.Printf("Timeout delegation ID%v, sending to local elevator: %v", order.ID, order)

				manager.orders.update(order)
				syncOrderWithOtherElevators(order, &manager)
				setHallLight(order.Dir, order.Floor, true)
			}

		case order := <-manager.orderCompleteTimeoutChannel:

			orderSaved, orderExists := manager.orders.getOrder(order.OwnerID, order.ID)

			if orderExists && orderSaved.State != msg.Completed {
				manager.logger.Printf("Order timeout ID%v: %#v", order.ID, order)
				selfServeHallOrder(order, &manager)
			}

		case peerUpdate := <-manager.peerUpdateChannel:
			for _, nodeid := range peerUpdate.Lost {

				manager.logger.Printf("Node lost connection: %v", nodeid)
				orders := manager.orders.getOrderToID(nodeid)

				for _, order := range orders {

					if order.OwnerID == manager.id {
						redelegateOrder(order, &manager)

					} else if !utility.IsStringInSlice(order.OwnerID, peerUpdate.Peers) &&
						!utility.IsStringInSlice(order.DelegatedToID, peerUpdate.Peers) {

						redelegateOrder(order, &manager)
					}
				}
			}

			if len(peerUpdate.New) > 0 {
				manager.logger.Printf("New node(s) connected")
				orders := manager.orders.getOrderFromID(manager.id)
				for _, order := range orders {
					if order.State == msg.Serving {
						syncOrderWithOtherElevators(order, &manager)
					}
				}
			}
		}
	}
}

func selfServeHallOrder(order msg.HallOrder, manager *HallOrderManager) {
	orderToFSM := elevio.ButtonEvent{
		Floor:  order.Floor,
		Button: elevio.ButtonType(order.Dir)}

	manager.delegateToLocalElevator <- orderToFSM
}

func syncOrderWithOtherElevators(order msg.HallOrder, manager *HallOrderManager) {
	manager.orderSyncToNetwork <- order
	manager.logger.Printf("Sync order ID%v to net:%#v", order.ID, order)
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

		manager.requestElevatorCost <- msg.RequestCost{
			Order: msg.OrderStamped{
				OrderID: order.ID,
				Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}},
			RequestFrom: msg.HallOrderManager}

		order.Costs[manager.id] = <-manager.elevatorCost

		manager.orders.update(order)

		timer.SendWithDelay(config.ORDER_REPLY_TIME, manager.orderReplyTimeoutChannel, order.ID)
		timer.SendWithDelayHallOrder(config.ORDER_COMPLETION_TIMEOUT, manager.orderCompleteTimeoutChannel, order)

		orderToNet := msg.OrderStamped{
			OrderID: order.ID,
			Order:   msg.Order{Floor: order.Floor, Dir: order.Dir}}

		manager.requestToNetwork <- orderToNet
	}
}