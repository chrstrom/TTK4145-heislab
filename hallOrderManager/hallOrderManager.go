package hallOrderManager

import (
	"fmt"

	"../elevio"
	"../localOrderDelegation"
	"../network"
	"../timer"
)

func initializeManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	delegateHallOrderCh chan<- elevio.ButtonEvent,
	requestElevatorCostCh chan<- elevio.ButtonEvent,
	elevatorCostCh <-chan int,
	orderCompleteCh <-chan elevio.ButtonEvent,
	channels network.NetworkChannels) HallOrderManager {

	var manager HallOrderManager

	manager.id = id

	manager.orders = make(OrderMap)
	manager.orderIDCounter = 1

	manager.localRequestChannel = localRequestCh
	manager.requestToNetwork = channels.RequestToNetwork
	manager.delegateToNetwork = channels.DelegateOrderToNetwork
	manager.delegationConfirmToNetwork = channels.DelegationConfirmToNetwork
	manager.requestReplyFromNetwork = channels.RequestReplyFromNetwork
	manager.orderDelegationConfirmFromNetwork = channels.DelegationConfirmFromNetwork
	manager.delegationFromNetwork = channels.DelegateFromNetwork

	manager.delegateToLocalElevator = delegateHallOrderCh
	manager.elevatorCost = elevatorCostCh
	manager.requestElevatorCost = requestElevatorCostCh
	manager.orderComplete = orderCompleteCh

	manager.orderReplyTimeoutChannel = make(chan int)
	manager.orderDelegationTimeoutChannel = make(chan int)

	return manager
}

func OrderManager(
	id string,
	localRequestCh <-chan localOrderDelegation.LocalOrder,
	delegateHallOrderCh chan<- elevio.ButtonEvent,
	requestElevatorCostCh chan<- elevio.ButtonEvent,
	elevatorCostCh <-chan int,
	orderCompleteCh <-chan elevio.ButtonEvent,
	channels network.NetworkChannels) {

	manager := initializeManager(id, localRequestCh, delegateHallOrderCh, requestElevatorCostCh, elevatorCostCh, orderCompleteCh, channels)

	for {
		select {
		case request := <-manager.localRequestChannel:
			//Check if order already exits? Or is this better to do in localOrdermanager? Or allow duplicates
			order := Order{OwnerID: manager.id,
				ID:    manager.orderIDCounter,
				State: Received,
				Floor: request.Floor,
				Dir:   request.Dir}
			manager.orderIDCounter++
			order.costs = make(map[string]int)

			//get local elevator cost in some way
			manager.requestElevatorCost <- elevio.ButtonEvent{Floor: order.Floor, Button: elevio.ButtonType(order.Dir)}
			order.costs[manager.id] = <-manager.elevatorCost
			fmt.Printf("Cost:%v\n", order.costs[manager.id])
			//order.costs[manager.id] = rand.Intn(1000)

			manager.orders.update(order)
			//fmt.Printf("%v - local request received \n", order.ID)

			timer.SendWithDelay(orderReplyTime, manager.orderReplyTimeoutChannel, order.ID)

			orderToNet := network.NewRequest{OrderID: order.ID, Floor: order.Floor, Dir: order.Dir}
			manager.requestToNetwork <- orderToNet

		case orderComplete := <-manager.orderComplete:
			//Get orders at the same floor from ordermap
			//Send confirmation to network
			fmt.Printf("Order %v completed\n", orderComplete)

		case reply := <-manager.requestReplyFromNetwork:
			o, valid := manager.orders.getOrder(manager.id, reply.OrderID)
			if valid && o.State == Received {
				o.costs[reply.ID] = reply.Cost
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
					manager.delegateToLocalElevator <- elevio.ButtonEvent{Floor: o.Floor, Button: elevio.ButtonType(o.Dir)}

					o.State = Serving
				} else {
					//fmt.Printf("%v - delegate to %v  (%v replies) \n", orderID, id, len(o.costs))
					timer.SendWithDelay(orderDelegationTime, manager.orderDelegationTimeoutChannel, orderID)

					o.State = Delegate

					message := network.Delegation{ID: o.DelegatedToID, OrderID: orderID, Floor: o.Floor, Dir: o.Dir}
					manager.delegateToNetwork <- message
				}
				manager.orders.update(o)
			}

		case delegation := <-manager.delegationFromNetwork:
			order := Order{OwnerID: delegation.ID,
				ID:            delegation.OrderID,
				DelegatedToID: manager.id,
				State:         Serving,
				Floor:         delegation.Floor,
				Dir:           delegation.Dir}
			manager.orders.update(order)
			reply := network.DelegationConfirm{ID: order.OwnerID, OrderID: order.ID, Floor: order.Floor, Dir: order.Dir}
			manager.delegationConfirmToNetwork <- reply

		case confirm := <-manager.orderDelegationConfirmFromNetwork:
			o, valid := manager.orders.getOrder(manager.id, confirm.OrderID)
			if valid && o.State == Delegate {
				o.State = Serving
				//fmt.Printf("%v - delegation confirmed \n", confirm.OrderID)

				manager.orders.update(o)
			}

		case orderID := <-manager.orderDelegationTimeoutChannel:
			o, valid := manager.orders.getOrder(manager.id, orderID)
			if valid && o.State == Delegate {
				//Send order to local elevator
				o.DelegatedToID = manager.id
				o.State = Serving
				manager.delegateToLocalElevator <- elevio.ButtonEvent{Floor: o.Floor, Button: elevio.ButtonType(o.Dir)}

				//fmt.Printf("%v - delegation timedout! Sending to local elevator \n", orderID)

				manager.orders.update(o)
			}
		}
	}
}
