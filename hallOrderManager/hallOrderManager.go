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

	manager.orders = make(map[int]Order)
	manager.orderIDCounter = 1

	manager.localRequestChannel = localRequestCh
	manager.requestToNetwork = channels.RequestToNetwork
	manager.delegateToNetwork = channels.DelegateOrderToNetwork
	manager.requestReplyFromNetwork = channels.RequestReplyFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork

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
			order := Order{State: Received, Floor: request.Floor, Dir: request.Dir}
			order.costs = make(map[string]int)
			//get local elevator cost in some way
			order.costs[manager.id] = rand.Intn(1000)

			orderID := manager.orderIDCounter
			manager.orderIDCounter++
			manager.orders[orderID] = order
			fmt.Printf("%v - local request received \n", orderID)

			timer.SendWithDelay(orderReplyTime, manager.orderReplyTimeoutChannel, orderID)

			orderToNet := network.NewRequest{OrderID: orderID, Floor: order.Floor, Dir: order.Dir}
			manager.requestToNetwork <- orderToNet

		case reply := <-manager.requestReplyFromNetwork:
			if manager.orders[reply.OrderID].State == Received && isValidOrder(manager.orders, reply.OrderID, reply.Floor, reply.Dir) {
				manager.orders[reply.OrderID].costs[reply.ID] = reply.Cost
			}

		case confirm := <-manager.orderDelegationConfirmFromNetwork:
			if manager.orders[confirm.OrderID].State == Delegate && isValidOrderConfirm(manager.orders, confirm.OrderID, confirm.Floor, confirm.Dir, confirm.ID) {
				o := manager.orders[confirm.OrderID]
				o.State = Serving
				fmt.Printf("%v - delegation confirmed \n", confirm.OrderID)
				// Send to Order Storage
				manager.orders[confirm.OrderID] = o
			}

		case orderID := <-manager.orderReplyTimeoutChannel:
			if manager.orders[orderID].State == Received && isValidOrderID(manager.orders, orderID) {
				o := manager.orders[orderID]
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

					message := network.Delegation{ID: o.DelegatedToID, OrderID: orderID, Floor: o.Floor, Dir: o.Dir}
					manager.delegateToNetwork <- message
				}
				manager.orders[orderID] = o
			}

		case orderID := <-manager.orderDelegationTimeoutChannel:
			if manager.orders[orderID].State == Delegate && isValidOrderID(manager.orders, orderID) {
				//Send order to local elevator
				o := manager.orders[orderID]
				o.DelegatedToID = manager.id
				o.State = Serving

				fmt.Printf("%v - delegation timedout! Sending to local elevator \n", orderID)
				// Send to Order Storage
				manager.orders[orderID] = o
			}
		}
	}
}

func isValidOrder(orders map[int]Order, orderID, floor, dir int) bool {
	o, found := orders[orderID]
	if found && o.Floor == floor && o.Dir == dir {
		return true
	}
	return false
}

func isValidOrderConfirm(orders map[int]Order, orderID, floor, dir int, id string) bool {
	o, found := orders[orderID]
	if found && o.Floor == floor && o.Dir == dir && o.DelegatedToID == id {
		return true
	}
	return false
}

func isValidOrderID(orders map[int]Order, orderID int) bool {
	_, found := orders[orderID]
	return found
}

func getIDOfLowestCost(costs map[string]int) string {
	lowest := 100000000
	lowestID := ""

	for id, c := range costs {
		if c <= lowest {
			lowest = c
			lowestID = id
		}
	}
	return lowestID
}
