package hallOrderManager

import (
	io "../elevio"
)

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

func setHallLight(dir int, floor int, state bool) {
	io.SetButtonLamp(io.ButtonType(dir), floor, state)
}