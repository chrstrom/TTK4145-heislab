package hallOrderManager

import (
	"../elevio"
)

func getIDOfLowestCost(costs map[string]int, defaultID string) string {
	lowest := 100000000
	lowestID := ""

	for id, c := range costs {
		if c <= lowest {
			lowest = c
			lowestID = id
		}
	}

	if lowestID == "" {
		lowestID = defaultID
	}

	return lowestID
}

func setHallLight(dir int, floor int, state bool) {
	elevio.SetButtonLamp(elevio.ButtonType(dir), floor, state)
}
