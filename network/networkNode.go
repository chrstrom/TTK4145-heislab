package network

import (
	"os"

	"log"

	"../network/bcast"
	"../network/peers"
	msg "../orderTypes"
)

const duplicatesOfMessages = 3

type Node struct {
	id               string
	messageIDCounter int

	networkChannels msg.NetworkChannels

	// Network channels
	peerUpdateChannelRx                                          chan peers.PeerUpdate
	peerTxEnable                                                 chan bool
	newRequestChannelTx, newRequestChannelRx                     chan msg.NetworkOrder
	newReplyToRequestChannelTx, newReplyToRequestChannelRx       chan msg.NetworkOrder
	delegateOrderChannelTx, delegateOrderChannelRx               chan msg.NetworkOrder
	delegateOrderConfirmChannelTx, delegateOrderConfirmChannelRx chan msg.NetworkOrder
	orderCompleteChannelTx, orderCompleteChannelRx               chan msg.NetworkOrder
	orderSyncChannelTx, orderSyncChannelRx                       chan msg.NetworkHallOrder

	receivedMessages map[string][]int

	loggerOutgoing, loggerIncoming *log.Logger
}

func NetworkNode(id string, channels msg.NetworkChannels) {

	node := initializeNetworkNode(id, channels)

	for {
		select {

		// Channels from the hall order manager to the network
		case request := <-node.networkChannels.RequestToNetwork:

			newRequest := msg.NetworkOrder{
				SenderID:  node.id,
				MessageID: node.messageIDCounter,
				Order:     request}

			node.messageIDCounter++

			node.loggerOutgoing.Printf("Request ID%v: %#v", newRequest.Order.OrderID, newRequest)
			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestChannelTx <- newRequest
			}

		case reply := <-node.networkChannels.ReplyToRequestToNetwork:

			newReplyToRequest := msg.NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: reply.ID,
				Order:      reply}

			node.messageIDCounter++

			node.loggerOutgoing.Printf("Reply to request: %#v", newReplyToRequest)
			for i := 0; i < duplicatesOfMessages; i++ {
				node.newReplyToRequestChannelTx <- newReplyToRequest
			}

		case delegation := <-node.networkChannels.DelegateOrderToNetwork:

			orderToBeDelegated := msg.NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: delegation.ID,
				Order:      delegation}

			node.messageIDCounter++
			//fmt.Printf("Delegate to network: %#v \n", orderToBeDelegated)

			node.loggerOutgoing.Printf("Delegation ID%v: %#v", orderToBeDelegated.Order.OrderID, orderToBeDelegated)
			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderChannelTx <- orderToBeDelegated
			}

		case confirm := <-node.networkChannels.DelegationConfirmToNetwork:

			confirmationOfDelegation := msg.NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: confirm.ID,
				Order:      confirm}

			node.messageIDCounter++
			//fmt.Printf("Sending Confirmation %#v \n", confirmationOfDelegation)

			node.loggerOutgoing.Printf("Confirmation of delegation: %#v", confirmationOfDelegation)
			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderConfirmChannelTx <- confirmationOfDelegation
			}

		case complete := <-node.networkChannels.OrderCompleteToNetwork:

			orderCompleted := msg.NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: complete.ID,
				Order:      complete}

			node.messageIDCounter++

			node.loggerOutgoing.Printf("Order completed ID%v: %#v", complete.OrderID, complete)
			for i := 0; i < duplicatesOfMessages; i++ {
				node.orderCompleteChannelTx <- orderCompleted
			}

		case order := <-node.networkChannels.SyncOrderToNetwork:

			syncOrder := msg.NetworkHallOrder{
				SenderID:  node.id,
				MessageID: node.messageIDCounter,
				Order:     order}
			node.messageIDCounter++

			node.loggerOutgoing.Printf("Sync order ID%v: %#v", order.ID, order)
			for i := 0; i < duplicatesOfMessages; i++ {
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
				//fmt.Printf("%#v \n", request)
				node.loggerIncoming.Printf("Request: %#v", request)
				message := msg.OrderStamped{
					ID:      request.SenderID,
					OrderID: request.Order.OrderID,
					Order:   request.Order.Order}

				node.networkChannels.RequestFromNetwork <- message
			}

		case replyToRequest := <-node.newReplyToRequestChannelRx:
			if replyToRequest.ReceiverID == node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					replyToRequest.SenderID,
					replyToRequest.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					replyToRequest.SenderID,
					replyToRequest.MessageID)
				//fmt.Printf("%#v \n", replyToRequest)
				node.loggerIncoming.Printf("Reply to request ID%v: %#v", replyToRequest.Order.OrderID, replyToRequest)
				message := msg.OrderStamped{
					ID:      replyToRequest.SenderID,
					OrderID: replyToRequest.Order.OrderID,
					Order:   replyToRequest.Order.Order}

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
				//fmt.Printf("Recieved delegation: %#v \n", delegation)
				node.loggerIncoming.Printf("Delegation: %#v", delegation)
				message := msg.OrderStamped{
					ID:      delegation.SenderID,
					OrderID: delegation.Order.OrderID,
					Order:   delegation.Order.Order}

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
				//fmt.Printf("Recieved confirmation: %#v \n", confirmation)
				node.loggerIncoming.Printf("Confirmation of delegation ID%v: %#v", confirmation.Order.OrderID, confirmation)
				message := msg.OrderStamped{
					ID:      confirmation.SenderID,
					OrderID: confirmation.Order.OrderID,
					Order:   confirmation.Order.Order}

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

				//fmt.Printf("%#v \n", complete)
				node.loggerIncoming.Printf("Complete ID%v: %#v", complete.Order.OrderID, complete)
				// Send message on channel
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
