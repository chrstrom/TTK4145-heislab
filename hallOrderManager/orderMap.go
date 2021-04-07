package hallOrderManager

import (
	"fmt"
	"sort"
)

type OrderMap map[string]map[int]HallOrder

func (om OrderMap) update(order HallOrder) {
	_, ok := om[order.OwnerID]
	if !ok {
		om[order.OwnerID] = make(map[int]HallOrder)
	}
	om[order.OwnerID][order.ID] = order

	om.printOrderMap()
}

func (om OrderMap) getOrder(ownerID string, orderID int) (order HallOrder, found bool) {
	_, ok := om[ownerID]
	if ok {
		o, ok2 := om[ownerID][orderID]
		if ok2 {
			return o, true
		}
	}
	return HallOrder{}, false
}

func (om OrderMap) printOrderMap() {
	//fmt.Print("\033[H\033[2J") //Clear screen in Go console
	/*cmd := exec.Command("cmd", "/c", "cls") //Clear screen in windows cmd
	cmd.Stdout = os.Stdout
	cmd.Run()*/

	fmt.Println("********************************OrderMap********************************")
	nodeids := make([]string, 0, len(om))
	var orders [][]int
	i := 0
	for id, omap := range om {
		nodeids = append(nodeids, id)
		orders = append(orders, []int{})
		for oid, _ := range omap {
			orders[i] = append(orders[i], oid)
		}
		sort.Ints(orders[i])
		i++
	}
	sort.Strings(nodeids)
	i = 0
	for _, id := range nodeids {
		fmt.Printf("Node: %s \n", id)
		fmt.Println("Order id    State        Delegated to          Floor    Direction")

		for _, oid := range orders[i] {
			o := om[id][oid]
			state := ""
			switch o.State {
			case Received:
				state = "Received"
			case Delegate:
				state = "Delegate"
			case Serving:
				state = "serving"
			}
			fmt.Printf("%v           %v      %s     %v        %v \n", o.ID, state, o.DelegatedToID, o.Floor, o.Dir)
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
