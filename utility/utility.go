package utility

import "strconv"

func StringArray2BoolArray(s []string) []bool {
	var b []bool
	for _, v := range s {
		bb, err := strconv.ParseBool(v)
		if err == nil {
			b = append(b, bb)
		}
	}
	return b
}

func FindMostCommonElement(s []string) (element string, count int) {
	countMap := make(map[string]int)
	for _, v := range s {
		countMap[v]++

		if countMap[v] > count {
			count = countMap[v]
			element = v
		}
	}

	return element, count
}
