package hallOrderManager

import (
	"sort"

	"../elevio"
)

func getIDOfLowestCost(costs map[string]int, defaultID string) string {
	lowest := 100000000
	lowestIDs := []string{}

	for id, c := range costs {
		if c < lowest {
			lowest = c
			lowestIDs = []string{id}
		} else if c == lowest {
			lowestIDs = append(lowestIDs, id)
		}
	}

	if len(lowestIDs) == 0 {
		lowestIDs = []string{defaultID}
	}
	sort.Strings(lowestIDs)
	return lowestIDs[0]
}

func setHallLight(dir int, floor int, state bool) {
	elevio.SetButtonLamp(elevio.ButtonType(dir), floor, state)
}
