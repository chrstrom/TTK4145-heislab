package network

import (
	"fmt"
	"os"
	"sort"

	"../network/localip"
	"../network/peers"
	types "../orderTypes"
)

func GetNodeID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%v-%v", localIP, os.Getpid())

	return id
}

func CreateNetworkChannelStruct() types.NetworkChannels {
	var networkChannels types.NetworkChannels

	networkChannels.RequestToNetwork = make(chan types.OrderStamped, 10)
	networkChannels.DelegateOrderToNetwork = make(chan types.OrderStamped, 10)
	networkChannels.ReplyToRequestToNetwork = make(chan types.OrderStamped, 10)
	networkChannels.DelegationConfirmToNetwork = make(chan types.OrderStamped, 10)
	networkChannels.OrderCompleteToNetwork = make(chan types.OrderStamped, 10)
	networkChannels.SyncOrderToNetwork = make(chan types.HallOrder, 10)

	networkChannels.RequestFromNetwork = make(chan types.OrderStamped, 10)
	networkChannels.DelegateFromNetwork = make(chan types.OrderStamped, 10)
	networkChannels.ReplyToRequestFromNetwork = make(chan types.OrderStamped, 10)
	networkChannels.DelegationConfirmFromNetwork = make(chan types.OrderStamped, 10)
	networkChannels.OrderCompleteFromNetwork = make(chan types.OrderStamped, 10)
	networkChannels.SyncOrderFromNetwork = make(chan types.HallOrder, 10)
	networkChannels.PeerUpdate = make(chan peers.PeerUpdate, 10)

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
