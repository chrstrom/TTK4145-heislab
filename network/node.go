package network

import (
	"fmt"
	"os"
	"sort"

	"../network-go/bcast"
	"../network-go/localip"
	"../network-go/peers"
)

func GetNodeID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%v-%v", localIP, os.Getpid())

	return id
}

type Node struct {
	id               string
	messageIDCounter int

	networkChannels NetworkChannels

	// Network channels
	peerUpdateChannel                                            chan peers.PeerUpdate
	peerTxEnable                                                 chan bool
	newRequestChannelTx, newRequestChannelRx                     chan NewRequestNetworkMessage
	newRequestReplyChannelTx, newRequestReplyChannelRx           chan NewRequestReplyNetworkMessage
	delegateOrderChannelTx, delegateOrderChannelRx               chan DelegateOrderNetworkMessage
	delegateOrderConfirmChannelTx, delegateOrderConfirmChannelRx chan DelegateOrderConfirmNetworkMessage
	orderCompleteChannelTx, orderCompleteChannelRx               chan OrderCompleteNetworkMessage

	receivedMessages map[string][]int
}

const duplicatesOfMessages = 3

func initializeNetworkNode(id string, channels NetworkChannels) Node {

	var node Node

	node.networkChannels = channels

	node.id = id
	node.messageIDCounter = 1

	node.peerUpdateChannel = make(chan peers.PeerUpdate)
	node.peerTxEnable = make(chan bool)
	go peers.Transmitter(25372, node.id, node.peerTxEnable)
	go peers.Receiver(25372, node.peerUpdateChannel)

	node.newRequestChannelTx = make(chan NewRequestNetworkMessage)
	node.newRequestChannelRx = make(chan NewRequestNetworkMessage)
	node.newRequestReplyChannelTx = make(chan NewRequestReplyNetworkMessage)
	node.newRequestReplyChannelRx = make(chan NewRequestReplyNetworkMessage)
	node.delegateOrderChannelTx = make(chan DelegateOrderNetworkMessage)
	node.delegateOrderChannelRx = make(chan DelegateOrderNetworkMessage)
	node.delegateOrderConfirmChannelTx = make(chan DelegateOrderConfirmNetworkMessage)
	node.delegateOrderConfirmChannelRx = make(chan DelegateOrderConfirmNetworkMessage)
	node.orderCompleteChannelTx = make(chan OrderCompleteNetworkMessage)
	node.orderCompleteChannelRx = make(chan OrderCompleteNetworkMessage)

	go bcast.Transmitter(25373, node.newRequestChannelTx)
	go bcast.Receiver(25373, node.newRequestChannelRx)

	go bcast.Transmitter(25374, node.newRequestReplyChannelTx)
	go bcast.Receiver(25374, node.newRequestReplyChannelRx)

	go bcast.Transmitter(25375, node.delegateOrderChannelTx)
	go bcast.Receiver(25375, node.delegateOrderChannelRx)

	go bcast.Transmitter(25376, node.delegateOrderConfirmChannelTx)
	go bcast.Receiver(25376, node.delegateOrderConfirmChannelRx)

	go bcast.Transmitter(25377, node.orderCompleteChannelTx)
	go bcast.Receiver(25377, node.orderCompleteChannelRx)

	node.receivedMessages = make(map[string][]int)

	return node
}

func NetworkNode(id string, channels NetworkChannels) {

	node := initializeNetworkNode(id, channels)

	for {
		select {
		case request := <-node.networkChannels.RequestFromNetwork:

			message := NewRequestNetworkMessage{
				SenderID:  node.id,
				MessageID: node.messageIDCounter,
				Floor:     request.Floor,
				Direction: request.Dir,
				OrderID:   request.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestChannelTx <- message
			}

		case reply := <-node.networkChannels.RequestReplyToNetwork:

			message := NewRequestReplyNetworkMessage{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: reply.ID,
				Floor:      reply.Floor,
				Direction:  reply.Dir,
				OrderID:    reply.OrderID,
				Cost:       reply.Cost}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestReplyChannelTx <- message
			}

		case delegation := <-node.networkChannels.DelegateFromNetwork:

			message := DelegateOrderNetworkMessage{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: delegation.ID,
				Floor:      delegation.Floor,
				Direction:  delegation.Dir,
				OrderID:    delegation.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderChannelTx <- message
			}

		case confirm := <-node.networkChannels.DelegationConfirmToNetwork:

			message := DelegateOrderConfirmNetworkMessage{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: confirm.ID,
				Floor:      confirm.Floor,
				Direction:  confirm.Dir,
				OrderID:    confirm.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderConfirmChannelTx <- message
			}

		case complete := <-node.networkChannels.OrderCompleteFromNetwork:

			message := OrderCompleteNetworkMessage{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: complete.ID,
				Floor:      complete.Floor,
				Direction:  complete.Dir,
				OrderID:    complete.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.orderCompleteChannelTx <- message
			}

		case request := <-node.newRequestChannelRx:
			if request.SenderID != node.id && shouldThisMessageBeProcessed(node.receivedMessages, request.SenderID, request.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, request.SenderID, request.MessageID)
				//fmt.Printf("%#v \n", request)

				message := NewRequest{
					ID:      request.SenderID,
					OrderID: request.OrderID,
					Floor:   request.Floor,
					Dir:     request.Direction}

				node.networkChannels.RequestToNetwork <- message
			}

		case requestReply := <-node.newRequestReplyChannelRx:
			if requestReply.ReceiverID == node.id && shouldThisMessageBeProcessed(node.receivedMessages, requestReply.SenderID, requestReply.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, requestReply.SenderID, requestReply.MessageID)
				//fmt.Printf("%#v \n", requestReply)

				message := RequestReply{ID: requestReply.SenderID,
					OrderID: requestReply.OrderID,
					Floor:   requestReply.Floor,
					Dir:     requestReply.Direction,
					Cost:    requestReply.Cost}

				node.networkChannels.RequestReplyFromNetwork <- message
			}

		case delegation := <-node.delegateOrderChannelRx:
			if delegation.ReceiverID == node.id && shouldThisMessageBeProcessed(node.receivedMessages, delegation.SenderID, delegation.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, delegation.SenderID, delegation.MessageID)
				//fmt.Printf("%#v \n", delegation)

				message := Delegation{ID: delegation.SenderID,
					OrderID: delegation.OrderID,
					Floor:   delegation.Floor,
					Dir:     delegation.Direction}

				node.networkChannels.DelegateOrderToNetwork <- message
			}

		case confirmation := <-node.delegateOrderConfirmChannelRx:
			if confirmation.ReceiverID == node.id && shouldThisMessageBeProcessed(node.receivedMessages, confirmation.SenderID, confirmation.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, confirmation.SenderID, confirmation.MessageID)
				//fmt.Printf("%#v \n", confirmation)

				message := DelegationConfirm{ID: confirmation.SenderID,
					OrderID: confirmation.OrderID,
					Floor:   confirmation.Floor,
					Dir:     confirmation.Direction}

				node.networkChannels.DelegationConfirmFromNetwork <- message
			}

		case complete := <-node.orderCompleteChannelRx:
			if complete.ReceiverID == node.id && shouldThisMessageBeProcessed(node.receivedMessages, complete.SenderID, complete.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, complete.SenderID, complete.MessageID)
				fmt.Printf("%#v \n", complete)

				// Send message on channel
			}

		}
	}
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
