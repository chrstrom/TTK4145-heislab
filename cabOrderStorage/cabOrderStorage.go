package cabOrderStorage

import (
	"io/ioutil"
	"strconv"
	"strings"
)

const numFloors = 4 //use istead a global variable for number of floors
const fileDuplicates = 3

func StoreCabOrders(orders [numFloors]bool) {
	orderString := ""
	for _, v := range orders {
		orderString = orderString + " " + strconv.FormatBool(v)
	}
	orderString = orderString[1:]

	for i := 0; i < fileDuplicates; i++ {
		filename := "cabOrderStorage/orderBackup/hallorders" + strconv.Itoa(i) + ".txt"
		ioutil.WriteFile(filename, []byte(orderString), 0644)
	}
}

func LoadCabOrders() [numFloors]bool {
	var orders [numFloors]bool

	fileDataString, allFilesEqual := readBackupFiles()

	if allFilesEqual {
		ordersString := strings.Split(fileDataString[0], " ")
		o := stringArray2BoolArray(ordersString)
		if len(o) == numFloors {
			copy(orders[:], o)
		}

	}

	return orders
}

func readBackupFiles() ([fileDuplicates]string, bool) {
	allFilesSame := true
	var fileData [fileDuplicates]string
	for i := 0; i < fileDuplicates; i++ {
		filename := "cabOrderStorage/orderBackup/hallorders" + strconv.Itoa(i) + ".txt"
		r, _ := ioutil.ReadFile(filename)
		fileData[i] = string(r)

		if fileData[i] != fileData[0] {
			allFilesSame = false
		}
	}
	return fileData, allFilesSame
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
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
