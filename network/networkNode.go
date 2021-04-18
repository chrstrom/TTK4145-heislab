package network

import (
	"fmt"
	"os"

	"log"

	"../config"
	msg "../messageTypes"
	"../network/bcast"
	"../network/peers"
)

// This is the driver function for the network node
// and contains a for-select, thus should be called as a goroutine.
func NetworkNode(id string, fsmChannels msg.FSMChannels, channels msg.NetworkChannels) {

	node := initializeNetworkNode(id, channels)

	for {
		select {

		// Channels from the hall order manager to the network
		case request := <-node.networkChannels.RequestToNetwork:

			newRequest := networkOrderFromOrderStamped(request, node)

			node.messageIDCounter++

			node.loggerOutgoing.Printf("Request ID%v: %#v", newRequest.Order.OrderID, newRequest)
			for i := 0; i < config.N_MESSAGE_DUPLICATES; i++ {
				node.newRequestChannelTx <- newRequest
			}

		case reply := <-fsmChannels.ReplyToNetWork:

			newReplyToRequest := networkOrderFromOrderStamped(reply, node)
			node.messageIDCounter++
			fmt.Printf("Network recieved cost: %#v\n", newReplyToRequest.Order.Cost)

			node.loggerOutgoing.Printf("Reply to request: %#v", newReplyToRequest)
			for i := 0; i < config.N_MESSAGE_DUPLICATES; i++ {
				node.newReplyToRequestChannelTx <- newReplyToRequest
			}

		case delegation := <-node.networkChannels.DelegateOrderToNetwork:

			orderToBeDelegated := networkOrderFromOrderStamped(delegation, node)
			node.messageIDCounter++
			node.loggerOutgoing.Printf("Delegation ID%v: %#v", orderToBeDelegated.Order.OrderID, orderToBeDelegated)

			for i := 0; i < config.N_MESSAGE_DUPLICATES; i++ {
				node.delegateOrderChannelTx <- orderToBeDelegated
			}

		case confirm := <-node.networkChannels.DelegationConfirmToNetwork:

			confirmationOfDelegation := networkOrderFromOrderStamped(confirm, node)

			node.messageIDCounter++

			node.loggerOutgoing.Printf("Confirmation of delegation: %#v", confirmationOfDelegation)
			for i := 0; i < config.N_MESSAGE_DUPLICATES; i++ {
				node.delegateOrderConfirmChannelTx <- confirmationOfDelegation
			}

		case complete := <-node.networkChannels.OrderCompleteToNetwork:

			orderCompleted := networkOrderFromOrderStamped(complete, node)
			node.messageIDCounter++
			node.loggerOutgoing.Printf("Order completed ID%v: %#v", complete.OrderID, complete)

			for i := 0; i < config.N_MESSAGE_DUPLICATES; i++ {
				node.orderCompleteChannelTx <- orderCompleted
			}

		case order := <-node.networkChannels.SyncOrderToNetwork:

			syncOrder := msg.NetworkHallOrder{
				SenderID:  node.id,
				MessageID: node.messageIDCounter,
				Order:     order}
			node.messageIDCounter++

			node.loggerOutgoing.Printf("Sync order ID%v: %#v", order.ID, order)
			for i := 0; i < config.N_MESSAGE_DUPLICATES; i++ {
				node.orderSyncChannelTx <- syncOrder
			}

			// Channels from the network to the hall order manager
		case request := <-node.newRequestChannelRx:
			if request.SenderID != node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					request.SenderID,
					request.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					request.SenderID,
					request.MessageID)

				node.loggerIncoming.Printf("Request: %#v", request)
				message := orderStampedFromNetworkOrder(request)
				fsmChannels.RequestCost <- msg.RequestCost{Order: message, RequestFrom: msg.Network}
			}

		case replyToRequest := <-node.newReplyToRequestChannelRx:
			if replyToRequest.ReceiverID == node.id &&
				shouldThisMessageBeProcessed(node.receivedMessages,
					replyToRequest.SenderID,
					replyToRequest.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					replyToRequest.SenderID,
					replyToRequest.MessageID)

				node.loggerIncoming.Printf("Reply to request ID%v: %#v", replyToRequest.Order.OrderID, replyToRequest)
				message := orderStampedFromNetworkOrder(replyToRequest)
				node.networkChannels.ReplyToRequestFromNetwork <- message
			}

		case delegation := <-node.delegateOrderChannelRx:
			if delegation.ReceiverID == node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					delegation.SenderID,
					delegation.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					delegation.SenderID,
					delegation.MessageID)

				node.loggerIncoming.Printf("Delegation: %#v", delegation)
				message := orderStampedFromNetworkOrder(delegation)
				node.networkChannels.DelegateFromNetwork <- message
			}

		case confirmation := <-node.delegateOrderConfirmChannelRx:
			if confirmation.ReceiverID == node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					confirmation.SenderID,
					confirmation.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					confirmation.SenderID,
					confirmation.MessageID)

				node.loggerIncoming.Printf("Confirmation of delegation ID%v: %#v", confirmation.Order.OrderID, confirmation)
				message := orderStampedFromNetworkOrder(confirmation)
				node.networkChannels.DelegationConfirmFromNetwork <- message
			}

		case complete := <-node.orderCompleteChannelRx:
			if complete.ReceiverID == node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					complete.SenderID,
					complete.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					complete.SenderID,
					complete.MessageID)

				node.loggerIncoming.Printf("Complete ID%v: %#v", complete.Order.OrderID, complete)
			}

		case sync := <-node.orderSyncChannelRx:
			if sync.SenderID != node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					sync.SenderID,
					sync.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					sync.SenderID,
					sync.MessageID)

				node.loggerIncoming.Printf("Sync order: %#v", sync)
				node.networkChannels.SyncOrderFromNetwork <- sync.Order

			}

		case peerUpdate := <-node.peerUpdateChannelRx:
			node.networkChannels.PeerUpdate <- peerUpdate
			node.loggerIncoming.Printf("Peer update: %#v", peerUpdate)

		}
	}
}

func initializeNetworkNode(id string, channels msg.NetworkChannels) Node {

	var node Node

	node.networkChannels = channels

	node.id = id
	node.messageIDCounter = 1

	node.peerUpdateChannelRx = make(chan peers.PeerUpdate)
	node.peerTxEnable = make(chan bool)
	go peers.Transmitter(25372, node.id, node.peerTxEnable)
	go peers.Receiver(25372, node.peerUpdateChannelRx)

	node.newRequestChannelTx = make(chan msg.NetworkOrder)
	node.newRequestChannelRx = make(chan msg.NetworkOrder)

	node.newReplyToRequestChannelTx = make(chan msg.NetworkOrder)
	node.newReplyToRequestChannelRx = make(chan msg.NetworkOrder)

	node.delegateOrderChannelTx = make(chan msg.NetworkOrder)
	node.delegateOrderChannelRx = make(chan msg.NetworkOrder)

	node.delegateOrderConfirmChannelTx = make(chan msg.NetworkOrder)
	node.delegateOrderConfirmChannelRx = make(chan msg.NetworkOrder)

	node.orderCompleteChannelTx = make(chan msg.NetworkOrder)
	node.orderCompleteChannelRx = make(chan msg.NetworkOrder)

	node.orderSyncChannelTx = make(chan msg.NetworkHallOrder)
	node.orderSyncChannelRx = make(chan msg.NetworkHallOrder)

	go bcast.Transmitter(25373, node.newRequestChannelTx)
	go bcast.Receiver(25373, node.newRequestChannelRx)

	go bcast.Transmitter(25374, node.newReplyToRequestChannelTx)
	go bcast.Receiver(25374, node.newReplyToRequestChannelRx)

	go bcast.Transmitter(25375, node.delegateOrderChannelTx)
	go bcast.Receiver(25375, node.delegateOrderChannelRx)

	go bcast.Transmitter(25376, node.delegateOrderConfirmChannelTx)
	go bcast.Receiver(25376, node.delegateOrderConfirmChannelRx)

	go bcast.Transmitter(25377, node.orderCompleteChannelTx)
	go bcast.Receiver(25377, node.orderCompleteChannelRx)

	go bcast.Transmitter(25378, node.orderSyncChannelTx)
	go bcast.Receiver(25378, node.orderSyncChannelRx)

	node.receivedMessages = make(map[string][]int)

	filepath := "log/" + node.id + "-network.log"
	file, _ := os.Create(filepath)
	node.loggerOutgoing = log.New(file, "Sending: ", log.Ltime|log.Lmicroseconds)
	node.loggerIncoming = log.New(file, "       --- Receiving: ", log.Ltime|log.Lmicroseconds)

	return node
}
