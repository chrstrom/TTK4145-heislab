package orderDelegation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/TTK4145-Students-2021/project-gruppe80/localOrderDelegation"
	"github.com/TTK4145-Students-2021/project-gruppe80/network"
	"github.com/TTK4145-Students-2021/project-gruppe80/timer"
)

type OrderStateType int

const (
	Received OrderStateType = iota
	Delegate
	Serving
)
const orderReplyTime = time.Millisecond * 50
const orderDelegationTime = time.Millisecond * 50

type Order struct {
	State         OrderStateType
	Floor, Dir    int
	costs         map[string]int
	DelegatedToID string
}

type OrderDelegator struct {
	id string

	orders         map[int]Order
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork  chan<- network.NewRequest
	delegateToNetwork chan<- network.Delegation

	requestReplyFromNetwork           <-chan network.RequestReply
	orderDelegationConfirmFromNetwork <-chan network.DelegationConfirm

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
}

func (delegator *OrderDelegator) Init(id string, localRequestCh <-chan localOrderDelegation.LocalOrder, requestToNetCh chan<- network.NewRequest,
	delegateToNetCh chan<- network.Delegation, requestReplyFromNetCh <-chan network.RequestReply, orderDelegationConfirmFromNetCh <-chan network.DelegationConfirm) {
	delegator.id = id

	delegator.orders = make(map[int]Order)
	delegator.orderIDCounter = 1

	delegator.localRequestChannel = localRequestCh
	delegator.requestToNetwork = requestToNetCh
	delegator.delegateToNetwork = delegateToNetCh
	delegator.requestReplyFromNetwork = requestReplyFromNetCh
	delegator.orderDelegationConfirmFromNetwork = orderDelegationConfirmFromNetCh

	delegator.orderReplyTimeoutChannel = make(chan int)
	delegator.orderDelegationTimeoutChannel = make(chan int)
}

func (delegator *OrderDelegator) OrderDelegation() {
	for {
		select {
		case request := <-delegator.localRequestChannel:
			//Check if order already exits? Or is this better to do in localOrderDelegator? Or allow duplicates
			order := Order{State: Received, Floor: request.Floor, Dir: request.Dir}
			order.costs = make(map[string]int)
			//get local elevator cost in some way
			order.costs[delegator.id] = rand.Intn(1000)

			orderID := delegator.orderIDCounter
			delegator.orderIDCounter++
			delegator.orders[orderID] = order
			fmt.Printf("%v - local request received \n", orderID)

			timer.SendWithDelay(orderReplyTime, delegator.orderReplyTimeoutChannel, orderID)

			orderToNet := network.NewRequest{OrderID: orderID, Floor: order.Floor, Dir: order.Dir}
			delegator.requestToNetwork <- orderToNet

		case reply := <-delegator.requestReplyFromNetwork:
			if isValidOrder(delegator.orders, reply.OrderID, reply.Floor, reply.Dir) {
				delegator.orders[reply.OrderID].costs[reply.ID] = reply.Cost
			}

		case confirm := <-delegator.orderDelegationConfirmFromNetwork:
			if isValidOrderConfirm(delegator.orders, confirm.OrderID, confirm.Floor, confirm.Dir, confirm.ID) && delegator.orders[confirm.OrderID].State == Delegate {
				o := delegator.orders[confirm.OrderID]
				o.State = Serving
				fmt.Printf("%v - delegation confirmed \n", confirm.OrderID)
				// Send to Order Storage
				delegator.orders[confirm.OrderID] = o
			}

		case orderID := <-delegator.orderReplyTimeoutChannel:
			if isValidOrderID(delegator.orders, orderID) && delegator.orders[orderID].State == Received {
				o := delegator.orders[orderID]
				id := getIDOfLowestCost(o.costs)
				if id == "" {
					id = delegator.id
				}
				o.DelegatedToID = id
				o.State = Delegate
				delegator.orders[orderID] = o

				if id == delegator.id {
					//send order to local elevator
					fmt.Printf("%v - delegate to local elevator (%v replies) \n", orderID, len(o.costs))
				} else {
					fmt.Printf("%v - delegate to %v  (%v replies) \n", orderID, id, len(o.costs))
					timer.SendWithDelay(orderDelegationTime, delegator.orderDelegationTimeoutChannel, orderID)
					message := network.Delegation{ID: o.DelegatedToID, OrderID: orderID, Floor: o.Floor, Dir: o.Dir}
					delegator.delegateToNetwork <- message
				}
			}

		case orderID := <-delegator.orderDelegationTimeoutChannel:
			if isValidOrderID(delegator.orders, orderID) && delegator.orders[orderID].State == Delegate {
				//Send order to local elevator
				o := delegator.orders[orderID]
				o.DelegatedToID = delegator.id
				o.State = Serving

				fmt.Printf("%v - delegation timedout! Sending to local elevator \n", orderID)
				// Send to Order Storage
				delegator.orders[orderID] = o
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
