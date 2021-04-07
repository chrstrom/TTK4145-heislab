package hallOrderManager

import (
	"fmt"
	"math/rand"

	"../localOrderDelegation"
	"../network"
	"../timer"
)

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

	manager.requestReplyFromNetwork = channels.RequestReplyFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork
	manager.orderSyncFromNetwork = channels.SyncOrderFromNetwork

	manager.orderReplyTimeoutChannel = make(chan int)
	manager.orderDelegationTimeoutChannel = make(chan int)

	return manager
}

func OrderManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	channels network.NetworkChannels) {

	manager := initializeManager(id, localRequestCh, channels)

	for {
		select {
		case request := <-manager.localRequestChannel:
			//Check if order already exits? Or is this better to do in localOrdermanager? Or allow duplicates
			order := HallOrder{
				ID: manager.orderIDCounter,
				OwnerID: manager.id,
				State:   Received,
				Floor: request.Floor,
				Dir: request.Dir}

			manager.orderIDCounter++
			order.costs = make(map[string]int)
			//get local elevator cost in some way
			order.costs[manager.id] = rand.Intn(1000)

			manager.orders.update(order)
			fmt.Printf("%v - local request received \n", order.ID)

			timer.SendWithDelay(orderReplyTime, manager.orderReplyTimeoutChannel, order.ID)

			orderToNet := network.OrderStamped{
				OrderID: order.ID,
				Order:   network.Order{Floor: order.Floor, Dir: order.Dir}}

			manager.requestToNetwork <- orderToNet

		case reply := <-manager.requestReplyFromNetwork:
			o, valid := manager.orders.getOrder(manager.id, reply.OrderID)
			if valid && o.State == Received {
				o.costs[reply.ID] = reply.Order.Cost
			}

		case confirm := <-manager.orderDelegationConfirmFromNetwork:
			o, valid := manager.orders.getOrder(manager.id, confirm.OrderID)
			if valid && o.State == Delegate {
				o.State = Serving
				fmt.Printf("%v - delegation confirmed \n", confirm.OrderID)

				manager.orders.update(o)

				//Let the elevators on the network know that this local elevator has taken an order
				orderToNet := network.OrderSync{OrderID: o.ID, Floor: o.Floor, Dir: o.Dir}
				manager.orderSyncToNetwork <- orderToNet
			}

		case sync := <-manager.orderSyncFromNetwork:
			if sync.ID != manager.id {

				order := HallOrder{}
				manager.orders.update(order)
			}

		case orderID := <-manager.orderReplyTimeoutChannel:
			o, valid := manager.orders.getOrder(manager.id, orderID)
			if valid && o.State == Received {
				id := getIDOfLowestCost(o.costs)
				if id == "" {
					id = manager.id
				}
				o.DelegatedToID = id

				if id == manager.id {
					//send order to local elevator
					fmt.Printf("%v - delegate to local elevator (%v replies) \n", orderID, len(o.costs))

					o.State = Serving
				} else {
					fmt.Printf("%v - delegate to %v  (%v replies) \n", orderID, id, len(o.costs))
					timer.SendWithDelay(orderDelegationTime, manager.orderDelegationTimeoutChannel, orderID)

					o.State = Delegate

					message := network.OrderStamped{
						ID:      o.DelegatedToID,
						OrderID: orderID,
						Order:   network.Order{Floor: o.Floor, Dir: o.Dir}}

					manager.delegateToNetwork <- message
				}
				manager.orders.update(o)
			}

		case orderID := <-manager.orderDelegationTimeoutChannel:
			o, valid := manager.orders.getOrder(manager.id, orderID)
			if valid && o.State == Delegate {
				//Send order to local elevator
				o.DelegatedToID = manager.id
				o.State = Serving

				fmt.Printf("%v - delegation timedout! Sending to local elevator \n", orderID)

				manager.orders.update(o)
			}
		}
	}
}
