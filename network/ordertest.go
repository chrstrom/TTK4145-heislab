package network

import (
	"fmt"
	"math/rand"
	"time"
)

type OrderStateType int

const (
	Recived OrderStateType = iota
	Delegate
	Serving
)
const ORDER_REPLY_TIME = time.Millisecond * 50
const ORDER_DELEGATION_TIME = ORDER_REPLY_TIME + time.Millisecond*50

type OrderCost struct {
	id   string
	cost int
}

type Order struct {
	state            OrderStateType
	floor, direction int
	costs            []OrderCost
	recivedTime      time.Time
}

type NewLocalOrder struct {
	floor, direction int
	cost             int
}

func OrderTest() {
	var node ElevatorNode
	node.Init()
	rand.Seed(time.Now().UTC().UnixNano())

	localOrderReciveChannel := make(chan NewLocalOrder)
	go makeRandomLocalRequests(localOrderReciveChannel)
	var orders []Order

	for {
		select {
		case newLocalRequest := <-localOrderReciveChannel:
			fmt.Printf("%v - New local order to floor\n", newLocalRequest.floor)
			newOrder := Order{state: Recived, floor: newLocalRequest.floor, direction: newLocalRequest.direction, recivedTime: time.Now()}
			localOrderCost := OrderCost{id: node.id, cost: newLocalRequest.cost}
			newOrder.costs = append(newOrder.costs, localOrderCost)
			orders = append(orders, newOrder)

			node.SendNewRequest(newLocalRequest.floor, newLocalRequest.direction)

		case newNetworkRequest := <-node.newRequestChannelRx:
			if newNetworkRequest.SenderID != node.id {
				fmt.Printf("    net %v - New Network request to floor \n", newNetworkRequest.Floor)
				cost := rand.Intn(1000) //random cost for testing
				node.SendNewReqestReply(newNetworkRequest.Floor, newNetworkRequest.Direction, cost)
			}

		case newRequestReply := <-node.newRequestReplyChannelRx:
			for i := range orders {
				if orders[i].state == Recived && orders[i].floor == newRequestReply.Floor && orders[i].direction == newRequestReply.Direction {
					cost := OrderCost{id: newRequestReply.SenderID, cost: newRequestReply.Cost}
					orders[i].costs = append(orders[i].costs, cost)
				}
			}

		case delegation := <-node.delegateOrderChannelRx:
			if delegation.ReceiverID == node.id {
				fmt.Printf("    net %v - Recived delegation to floor \n", delegation.Floor)
				//Send order to local elevator
				node.SendDelegateOrderConfirm(delegation.SenderID, delegation.Floor, delegation.Direction)
			}

		case confirmation := <-node.delegateOrderConfirmChannelRx:
			if confirmation.ReceiverID == node.id {
				for i := range orders {
					if orders[i].state == Delegate && orders[i].floor == confirmation.Floor && orders[i].direction == confirmation.Direction {
						fmt.Printf("%v - Order successfully delegated to remote elevator\n", confirmation.Floor)
						orders[i].state = Serving
					}
				}
			}

		default:
			// Tror dette er en dårlig løsning, kanskje bedre om det sendes en melding på en kanal når tiden
			// går ut. Da vil denne funksjonaliteten legges inn i egne case, istenden for i default
			for i := range orders {
				switch orders[i].state {
				case Recived:
					if time.Now().Sub(orders[i].recivedTime) >= ORDER_REPLY_TIME {
						orders[i].delegateOrder(&node)
					}

				case Delegate:
					if time.Now().Sub(orders[i].recivedTime) >= ORDER_DELEGATION_TIME {
						fmt.Printf("%v - Did not get delegation confirmation, delegate to order to me \n", orders[i].floor)
						// Send order to local elevator
						orders[i].state = Serving
					}
				}
			}
			time.Sleep(time.Millisecond * 10)
		}

	}
}

func (order *Order) delegateOrder(node *ElevatorNode) {
	lowest := 10000000000000
	var id string
	for _, c := range order.costs {
		if c.cost < lowest {
			lowest = c.cost
			id = c.id
		}
	}

	if id == node.id {
		fmt.Printf("%v - Delegate this order to me \n", order.floor)
		//Send order to local elevator
		order.state = Serving
	} else {
		node.SendDelegateOrder(id, order.floor, order.direction)
		order.state = Delegate
	}
}

func makeRandomLocalRequests(localOrderChannel chan<- NewLocalOrder) {
	order := NewLocalOrder{floor: rand.Intn(1000), direction: 0, cost: rand.Intn(1000)}
	for {
		localOrderChannel <- order
		//use floor as message id for this test
		order.floor++
		order.cost = rand.Intn(1000)
		time.Sleep(10 * time.Second)
	}
}
