package network

import (
	"fmt"

	"../network-go/bcast"
	"../network-go/peers"
)

const duplicatesOfMessages = 3

type Node struct {
	id               string
	messageIDCounter int

	networkChannels NetworkChannels

	// Network channels
	peerUpdateChannel                                            chan peers.PeerUpdate
	peerTxEnable                                                 chan bool
	newRequestChannelTx, newRequestChannelRx                     chan NetworkOrder
	newRequestReplyChannelTx, newRequestReplyChannelRx           chan NetworkOrder
	delegateOrderChannelTx, delegateOrderChannelRx               chan NetworkOrder
	delegateOrderConfirmChannelTx, delegateOrderConfirmChannelRx chan NetworkOrder
	orderCompleteChannelTx, orderCompleteChannelRx               chan NetworkOrder
	orderSyncChannelTx, orderSyncChannelRx                       chan OrderSyncNetworkMessage

	receivedMessages map[string][]int
}

func NetworkNode(id string, channels NetworkChannels) {

	node := initializeNetworkNode(id, channels)

	for {
		select {

		// Channels from the hall order manager to the network
		case request := <-node.networkChannels.RequestToNetwork:

			newRequest := NetworkOrder{
				SenderID:  node.id,
				MessageID: node.messageIDCounter,
				Floor:     request.Floor,
				Direction: request.Dir,
				OrderID:   request.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestChannelTx <- newRequest
			}

		case reply := <-node.networkChannels.RequestReplyToNetwork:

			newReplyToRequest := NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: reply.ID,
				Floor:      reply.Floor,
				Direction:  reply.Dir,
				OrderID:    reply.OrderID,
				Cost:       reply.Cost}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestReplyChannelTx <- newReplyToRequest
			}

		case delegation := <-node.networkChannels.DelegateOrderToNetwork:

			orderToBeDelegated := NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: delegation.ID,
				Floor:      delegation.Floor,
				Direction:  delegation.Dir,
				OrderID:    delegation.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderChannelTx <- orderToBeDelegated
			}

		case confirm := <-node.networkChannels.DelegationConfirmToNetwork:

			confirmationOfDelegation := NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: confirm.ID,
				Floor:      confirm.Floor,
				Direction:  confirm.Dir,
				OrderID:    confirm.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderConfirmChannelTx <- confirmationOfDelegation 
			}

		case complete := <-node.networkChannels.OrderCompleteToNetwork:

			orderCompleted := NetworkOrder{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: complete.ID,
				Floor:      complete.Floor,
				Direction:  complete.Dir,
				OrderID:    complete.OrderID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.orderCompleteChannelTx <- orderCompleted
			}

		case order := <-node.networkChannels.SyncOrderToNetwork:
			message := OrderSyncNetworkMessage{
				SenderID:   node.id,
				MessageID:  node.messageIDCounter,
				ReceiverID: order.ID}

			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.orderSyncChannelTx <- message
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

				message := Order{
					ID:      request.SenderID,
					OrderID: request.OrderID,
					Floor:   request.Floor,
					Dir:     request.Direction}

				node.networkChannels.RequestFromNetwork <- message
			}

		case replyToRequest := <-node.newRequestReplyChannelRx:
			if replyToRequest.ReceiverID == node.id &&
				shouldThisMessageBeProcessed(
					node.receivedMessages,
					replyToRequest.SenderID,
					replyToRequest.MessageID) {

				addMessageIDToReceivedMessageMap(
					node.receivedMessages,
					replyToRequest.SenderID,
					replyToRequest.MessageID)
				//fmt.Printf("%#v \n", requestReply)

				message := Order{ID: replyToRequest.SenderID,
					OrderID: replyToRequest.OrderID,
					Floor:   replyToRequest.Floor,
					Dir:     replyToRequest.Direction,
					Cost:    replyToRequest.Cost}

				node.networkChannels.RequestReplyFromNetwork <- message
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
				//fmt.Printf("%#v \n", delegation)

				message := Order{ID: delegation.SenderID,
					OrderID: delegation.OrderID,
					Floor:   delegation.Floor,
					Dir:     delegation.Direction}

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
				//fmt.Printf("%#v \n", confirmation)

				message := Order{ID: confirmation.SenderID,
					OrderID: confirmation.OrderID,
					Floor:   confirmation.Floor,
					Dir:     confirmation.Direction}

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

				fmt.Printf("%#v \n", complete)

				// Send message on channel
			}

		case _ = <-node.orderSyncChannelRx:
			message := OrderSync{}

			node.networkChannels.SyncOrderFromNetwork <- message

		}
	}
}

func initializeNetworkNode(id string, channels NetworkChannels) Node {

	var node Node

	node.networkChannels = channels

	node.id = id
	node.messageIDCounter = 1

	node.peerUpdateChannel = make(chan peers.PeerUpdate)
	node.peerTxEnable = make(chan bool)
	go peers.Transmitter(25372, node.id, node.peerTxEnable)
	go peers.Receiver(25372, node.peerUpdateChannel)

	node.newRequestChannelTx = make(chan NetworkOrder)
	node.newRequestChannelRx = make(chan NetworkOrder)

	node.newRequestReplyChannelTx = make(chan NetworkOrder)
	node.newRequestReplyChannelRx = make(chan NetworkOrder)

	node.delegateOrderChannelTx = make(chan NetworkOrder)
	node.delegateOrderChannelRx = make(chan NetworkOrder)

	node.delegateOrderConfirmChannelTx = make(chan NetworkOrder)
	node.delegateOrderConfirmChannelRx = make(chan NetworkOrder)

	node.orderCompleteChannelTx = make(chan NetworkOrder)
	node.orderCompleteChannelRx = make(chan NetworkOrder)

	node.orderSyncChannelTx = make(chan OrderSyncNetworkMessage)
	node.orderSyncChannelRx = make(chan OrderSyncNetworkMessage)

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

	go bcast.Transmitter(25378, node.orderSyncChannelTx)
	go bcast.Receiver(25378, node.orderSyncChannelRx)

	node.receivedMessages = make(map[string][]int)

	return node
}
