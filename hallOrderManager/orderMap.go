package hallOrderManager

import (
	"fmt"
	"sort"

	msg "../messageTypes"
)

type OrderMap map[string]map[int]msg.HallOrder

func (orderMap OrderMap) update(order msg.HallOrder) {
	_, ok := orderMap[order.OwnerID]
	if !ok {
		orderMap[order.OwnerID] = make(map[int]msg.HallOrder)
	}
	orderMap[order.OwnerID][order.ID] = order

	orderMap.print()
}

func (orderMap OrderMap) getOrder(ownerID string, orderID int) (order msg.HallOrder, found bool) {
	_, ok := orderMap[ownerID]
	if ok {
		order, ok2 := orderMap[ownerID][orderID]
		if ok2 {
			return order, true
		}
	}
	return msg.HallOrder{}, false
}

func (orderMap OrderMap) getOrderFromID(ownerID string) []msg.HallOrder {
	var orders []msg.HallOrder
	_, ok := orderMap[ownerID]
	if ok {
		for _, order := range orderMap[ownerID] {
			orders = append(orders, order)
		}
	}
	return orders
}

func (orderMap OrderMap) getOrderToID(delegatedID string) []msg.HallOrder {
	var orders []msg.HallOrder
	for _, node := range orderMap {
		for _, order := range node {
			if order.DelegatedToID == delegatedID {
				orders = append(orders, order)
			}
		}
	}
	return orders
}

func (orderMap OrderMap) getOrdersToFloorWithDir(floor, dir int) []msg.HallOrder {

	orders := make([]msg.HallOrder, 0)

	for _, node := range orderMap {
		for _, order := range node {
			if order.Dir == dir && order.Floor == floor {
				orders = append(orders, order)
			}
		}
	}
	return orders
}

func (orderMap OrderMap) print() {
	// Please beware that this function is UGLY, but pretty printing usually is,
	// so fuck it B)

	fmt.Println("********************************OrderMap********************************")

	//Reformate the order map into one integer slice per node
	nodes := []struct {
		nodeid  string
		orderid []int
	}{}
	i := 0
	for id, singleNodeOrders := range orderMap {
		nodes = append(nodes, struct {
			nodeid  string
			orderid []int
		}{id, []int{}})

		for oid := range singleNodeOrders {
			nodes[i].orderid = append(nodes[i].orderid, oid)
		}
		sort.Ints(nodes[i].orderid)
		i++
	}
	//Sort nodes on id
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].nodeid < nodes[j].nodeid })

	//Print orders from each node
	for _, node := range nodes {
		fmt.Printf("Node: %s \n", node.nodeid)
		fmt.Println("Order id    State        Delegated to          Floor    Direction")

		for _, oid := range node.orderid {
			o := orderMap[node.nodeid][oid]
			state := ""
			switch o.State {
			case msg.Received:
				state = "Received"
			case msg.Delegate:
				state = "Delegate"
			case msg.Serving:
				state = "Serving"
			case msg.Completed:
				state = "Completed"
			}
			fmt.Printf("%-11v %-12v %-21s %-8v %v \n", o.ID, state, o.DelegatedToID, o.Floor, o.Dir)
		}
		fmt.Printf("\n\n")
	}
}
