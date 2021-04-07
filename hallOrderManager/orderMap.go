package hallOrderManager

type OrderMap map[string]map[int]HallOrder

func (om OrderMap) update(order HallOrder) {
	_, ok := om[order.OwnerID]
	if !ok {
		om[order.OwnerID] = make(map[int]HallOrder)
	}
	om[order.OwnerID][order.ID] = order
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
