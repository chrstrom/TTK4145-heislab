package network

import (
	"fmt"
	"os"
	"sort"

	"../config"
	msg "../messageTypes"
	"../network/localip"
	"../network/peers"
)

func GetNodeID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%v-%v", localIP, os.Getpid())

	return id
}

func CreateNetworkChannelStruct() msg.NetworkChannels {
	var networkChannels msg.NetworkChannels

	const bufferSize = config.CHANNEL_BUFFER_SIZE

	networkChannels.RequestToNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.DelegateOrderToNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.DelegationConfirmToNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.OrderCompleteToNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.SyncOrderToNetwork = make(chan msg.HallOrder, bufferSize)

	networkChannels.DelegateFromNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.ReplyToRequestFromNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.DelegationConfirmFromNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.OrderCompleteFromNetwork = make(chan msg.OrderStamped, bufferSize)
	networkChannels.SyncOrderFromNetwork = make(chan msg.HallOrder, bufferSize)
	networkChannels.PeerUpdate = make(chan peers.PeerUpdate, bufferSize)

	return networkChannels
}

func networkOrderFromOrderStamped(order msg.OrderStamped, node Node) msg.NetworkOrder {
	return msg.NetworkOrder{
		SenderID:   node.id,
		MessageID:  node.messageIDCounter,
		ReceiverID: order.ID,
		Order:      order}
}

func orderStampedFromNetworkOrder(order msg.NetworkOrder) msg.OrderStamped {
	return msg.OrderStamped{
		ID:      order.SenderID,
		OrderID: order.Order.OrderID,
		Floor:   order.Order.Floor,
		Dir:     order.Order.Dir,
		Cost:    order.Order.Cost}
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
