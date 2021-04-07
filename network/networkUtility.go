package network

import (
	"fmt"
	"os"
	"sort"

	"../network/localip"
	"../hallOrderManager"
)

func GetNodeID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%v-%v", localIP, os.Getpid())

	return id
}

func CreateNetworkChannelStruct() NetworkChannels {
	var networkChannels NetworkChannels

	networkChannels.RequestToNetwork 				= make(chan OrderStamped)
	networkChannels.DelegateOrderToNetwork 			= make(chan OrderStamped)
	networkChannels.RequestReplyToNetwork 			= make(chan OrderStamped)
	networkChannels.DelegationConfirmToNetwork 		= make(chan OrderStamped)
	networkChannels.OrderCompleteToNetwork 			= make(chan OrderStamped)
	networkChannels.SyncOrderToNetwork				= make(chan hallOrderManager.HallOrder)

	networkChannels.RequestFromNetwork 				= make(chan OrderStamped)
	networkChannels.DelegateFromNetwork 			= make(chan	OrderStamped)
	networkChannels.RequestReplyFromNetwork 		= make(chan OrderStamped)
	networkChannels.DelegationConfirmFromNetwork 	= make(chan OrderStamped)
	networkChannels.OrderCompleteFromNetwork 		= make(chan OrderStamped)
	networkChannels.SyncOrderFromNetwork			= make(chan hallOrderManager.HallOrder)

	return networkChannels
}

func shouldThisMessageBeProcessed(receivedMessages map[string][]int, senderID string, messageID int) bool {
	process := true
	s, exists := receivedMessages[senderID]
	if exists {
		i := sort.SearchInts(receivedMessages[senderID], messageID)
		if messageID < s[0] || (i < len(s) && s[i] == messageID) {
			process = false
		}
	}
	return process
}

func addMessageIDToReceivedMessageMap(receivedMessages map[string][]int, senderID string, messageID int) {
	_, exists := receivedMessages[senderID]
	if exists == false {
		addNodeToMessageMap(receivedMessages, senderID)
	}
	receivedMessages[senderID][0] = messageID
	sort.Ints(receivedMessages[senderID])
}

func addNodeToMessageMap(mm map[string][]int, nodeID string) {
	mm[nodeID] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}
