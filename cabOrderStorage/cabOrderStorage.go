package cabOrderStorage

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/TTK4145-Students-2021/project-gruppe80/utility"
)

const numFloors = 4 //use istead a global variable for number of floors
const fileDuplicates = 3
const backupPath = "orderBackup/"

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
		e, count := utility.FindMostCommonElement(fileData[:])
		equalFilesRequired := int(math.Ceil(fileDuplicates / 2.0))

		if count >= equalFilesRequired {
			ordersString = e
		}
	}

	split := strings.Split(ordersString, " ")
	o := utility.StringArray2BoolArray(split)
	if len(split) == numFloors && len(o) == numFloors {
		copy(orders[:], o)
	}

	return orders
}

func readBackupFiles() (fileData [fileDuplicates]string, allFilesSame bool) {
	allFilesSame = true
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
