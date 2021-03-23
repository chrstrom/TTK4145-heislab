package cabOrderStorage

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

const numFloors = 4 //use istead a global variable for number of floors
const fileDuplicates = 3
const backupPath = "cabOrderStorage/orderBackup/"

func StoreCabOrders(orders [numFloors]bool) {
	orderString := ""
	for _, v := range orders {
		orderString = orderString + " " + strconv.FormatBool(v)
	}
	orderString = orderString[1:]

	for i := 0; i < fileDuplicates; i++ {
		filename := backupPath + "hallorders" + strconv.Itoa(i) + ".txt"
		ioutil.WriteFile(filename, []byte(orderString), 0644)
	}
}

func LoadCabOrders() [numFloors]bool {
	var orders [numFloors]bool
	var ordersString string

	fileData, allFilesEqual := readBackupFiles()

	if allFilesEqual {
		ordersString = fileData[0]
	} else {
		e, count := findMostCommonElement(fileData)
		equalFilesRequired := int(math.Ceil(fileDuplicates / 2.0))

		if count >= equalFilesRequired {
			ordersString = e
		}
	}

	split := strings.Split(ordersString, " ")
	o := stringArray2BoolArray(split)
	if len(split) == numFloors && len(o) == numFloors {
		copy(orders[:], o)
	}

	return orders
}

func readBackupFiles() (fileData [fileDuplicates]string, allFilesSame bool) {
	for i := 0; i < fileDuplicates; i++ {
		filename := backupPath + "hallorders" + strconv.Itoa(i) + ".txt"
		r, _ := ioutil.ReadFile(filename)
		fileData[i] = string(r)

		if fileData[i] != fileData[0] {
			allFilesSame = false
		}
	}
	return fileData, allFilesSame
}

func stringArray2BoolArray(s []string) []bool {
	var b []bool
	for _, v := range s {
		bb, err := strconv.ParseBool(v)
		if err == nil {
			b = append(b, bb)
		}
	}
	return b
}

func findMostCommonElement(s [fileDuplicates]string) (element string, count int) {
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
