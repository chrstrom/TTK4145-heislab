package hallOrderManager

import (
	"fmt"
	"sort"

	msg "../orderTypes"
)

type OrderMap map[string]map[int]msg.HallOrder

func (om OrderMap) update(order msg.HallOrder) {
	_, ok := om[order.OwnerID]
	if !ok {
		om[order.OwnerID] = make(map[int]msg.HallOrder)
	}
	om[order.OwnerID][order.ID] = order

	om.printOrderMap()
}

func (om OrderMap) getOrder(ownerID string, orderID int) (order msg.HallOrder, found bool) {
	_, ok := om[ownerID]
	if ok {
		o, ok2 := om[ownerID][orderID]
		if ok2 {
			return o, true
		}
	}
	return msg.HallOrder{}, false
}

/*func (om OrderMap) getOrderFromIDtoID(ownerID, delegatedID string) []msg.HallOrder {
	var orders []msg.HallOrder
	_, ok := om[ownerID]
	if ok {
		for _, o := range om[ownerID] {
			if o.DelegatedToID == delegatedID {
				orders = append(orders, o)
			}
		}
	}
	return orders
}*/

func (om OrderMap) getOrderFromID(ownerID string) []msg.HallOrder {
	var orders []msg.HallOrder
	_, ok := om[ownerID]
	if ok {
		for _, o := range om[ownerID] {
			orders = append(orders, o)
		}
	}
	return orders
}

func (om OrderMap) getOrderToID(delegatedID string) []msg.HallOrder {
	var orders []msg.HallOrder
	for _, node := range om {
		for _, o := range node {
			if o.DelegatedToID == delegatedID {
				orders = append(orders, o)
			}
		}
	}
	return orders
}

func (om OrderMap) getOrdersToFloorWithDir(floor, dir int) []msg.HallOrder {

	orders := make([]msg.HallOrder, 0)

	for _, node := range om {
		for _, order := range node {
			if order.Dir == dir && order.Floor == floor {
				orders = append(orders, order)
			}
		}
	}
	return orders
}

func (om OrderMap) printOrderMap() {
	//fmt.Print("\033[H\033[2J") //Clear screen in Go console
	/*cmd := exec.Command("cmd", "/c", "cls") //Clear screen in windows cmd
	cmd.Stdout = os.Stdout
	cmd.Run()*/

	fmt.Println("********************************OrderMap********************************")
	orders := []struct {
		nodeid  string
		orderid []int
	}{}
	i := 0
	for id, omap := range om {
		orders = append(orders, struct {
			nodeid  string
			orderid []int
		}{id, []int{}})
		for oid, _ := range omap {
			orders[i].orderid = append(orders[i].orderid, oid)
		}
		sort.Ints(orders[i].orderid)
		i++
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].nodeid < orders[j].nodeid })
	i = 0
	for _, node := range orders {
		fmt.Printf("Node: %s \n", node.nodeid)
		fmt.Println("Order id    State        Delegated to          Floor    Direction")

		for _, oid := range node.orderid {
			o := om[node.nodeid][oid]
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
		i++
		fmt.Printf("\n\n")
	}
}

/*
func (om OrderMap) isValidOrder(order Order) bool {
	_, ok := om[order.OwnerID]
	if ok {
		o, ok2 := om[order.OwnerID][order.ID]
		if ok2 && o.ID == order.ID && o.Floor == order.Floor && o.Dir == order.Dir {
			return true
		}
	}
	return false
}*/
/*
func TestOrderMap() {
	om := make(OrderMap)
	o := Order{ID: 123, State: Received, Floor: 2, Dir: 1, OwnerID: "node1"}
	//o2 := Order{ID: 123, State: Received, Floor: 2, Dir: 1, OwnerID: "node1"}
	om.update(o)
	fmt.Println(om)
	fmt.Println(om.getOrder("node1", 123))
	o3, _ := om.getOrder("node1", 123)
	o3.Floor = 20
	om.update(o3)
	fmt.Println(om)
}*/
