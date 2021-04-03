package hallOrderManager

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
