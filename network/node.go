package network

import (
	"fmt"
	"os"
	"sort"

	"github.com/TTK4145-Students-2021/project-gruppe80/network-Go-modul/bcast"
	"github.com/TTK4145-Students-2021/project-gruppe80/network-Go-modul/localip"
	"github.com/TTK4145-Students-2021/project-gruppe80/network-Go-modul/peers"
)

type Node struct {
	id               string
	messageIDCounter int

	// Local channels
	requestLocalChannelIn           <-chan NewRequest
	delegateOrderLocalChannelIn     <-chan Delegation
	requestReplyLocalChannelIn      <-chan RequestReply
	delegationComfirmLocalChannelIn <-chan DelegationConfirm
	orderCompleteLocalChannelIn     <-chan OrderComplete

	requestLocalChannelOut           chan<- NewRequest
	delegateOrderLocalChannelOut     chan<- Delegation
	requestReplyLocalChannelOut      chan<- RequestReply
	delegationComfirmLocalChannelOut chan<- DelegationConfirm
	orderCompleteLocalChannelOut     chan<- OrderComplete

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

func (node *Node) Init(requestChIn <-chan NewRequest, delegationChIn <-chan Delegation, requestReplyChIn <-chan RequestReply,
	delegationComfirmChIn <-chan DelegationConfirm, orderCompleteChIn <-chan OrderComplete, requestChOut chan<- NewRequest,
	delegationChOut chan<- Delegation, requestReplyChOut chan<- RequestReply, delegationComfirmChOut chan<- DelegationConfirm,
	orderCompleteChOut chan<- OrderComplete) {

	node.requestLocalChannelIn = requestChIn
	node.delegateOrderLocalChannelIn = delegationChIn
	node.requestReplyLocalChannelIn = requestReplyChIn
	node.delegationComfirmLocalChannelIn = delegationComfirmChIn
	node.orderCompleteLocalChannelIn = orderCompleteChIn

	node.requestLocalChannelOut = requestChOut
	node.delegateOrderLocalChannelOut = delegationChOut
	node.requestReplyLocalChannelOut = requestReplyChOut
	node.delegationComfirmLocalChannelOut = delegationComfirmChOut
	node.orderCompleteLocalChannelOut = orderCompleteChOut

	localIP, err := localip.LocalIP()
	if err != nil {
		localIP = "DISCONNECTED"
	}
	node.id = fmt.Sprintf("%v-%v", localIP, os.Getpid())
	fmt.Printf("Init elevator network node with id:%v\n", node.id)

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
}

func (node *Node) NetworkNode() {
	for {
		select {
		case request := <-node.requestLocalChannelIn:
			message := NewRequestNetworkMessage{SenderID: node.id, MessageID: node.messageIDCounter, Floor: request.Floor, Direction: request.Dir, OrderID: request.OrderID}
			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestChannelTx <- message
			}

		case reply := <-node.requestReplyLocalChannelIn:
			message := NewRequestReplyNetworkMessage{SenderID: node.id, MessageID: node.messageIDCounter, ReceiverID: reply.ID, Floor: reply.Floor, Direction: reply.Dir, OrderID: reply.OrderID, Cost: reply.Cost}
			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.newRequestReplyChannelTx <- message
			}

		case delegation := <-node.delegateOrderLocalChannelIn:
			message := DelegateOrderNetworkMessage{SenderID: node.id, MessageID: node.messageIDCounter, ReceiverID: delegation.ID, Floor: delegation.Floor, Direction: delegation.Dir, OrderID: delegation.OrderID}
			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderChannelTx <- message
			}

		case confirm := <-node.delegationComfirmLocalChannelIn:
			message := DelegateOrderConfirmNetworkMessage{SenderID: node.id, MessageID: node.messageIDCounter, ReceiverID: confirm.ID, Floor: confirm.Floor, Direction: confirm.Dir, OrderID: confirm.OrderID}
			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.delegateOrderConfirmChannelTx <- message
			}

		case complete := <-node.orderCompleteLocalChannelIn:
			message := OrderCompleteNetworkMessage{SenderID: node.id, MessageID: node.messageIDCounter, ReceiverID: complete.ID, Floor: complete.Floor, Direction: complete.Dir, OrderID: complete.OrderID}
			node.messageIDCounter++

			for i := 0; i < duplicatesOfMessages; i++ {
				node.orderCompleteChannelTx <- message
			}

		case request := <-node.newRequestChannelRx:
			if shouldThisMessageBeProcessed(node.receivedMessages, request.SenderID, request.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, request.SenderID, request.MessageID)
				fmt.Printf("%#v \n", request)
			}

		case requestReply := <-node.newRequestReplyChannelRx:
			if shouldThisMessageBeProcessed(node.receivedMessages, requestReply.SenderID, requestReply.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, requestReply.SenderID, requestReply.MessageID)
				fmt.Printf("%#v \n", requestReply)
			}

		case delegation := <-node.delegateOrderChannelRx:
			if shouldThisMessageBeProcessed(node.receivedMessages, delegation.SenderID, delegation.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, delegation.SenderID, delegation.MessageID)
				fmt.Printf("%#v \n", delegation)
			}

		case confirmation := <-node.delegateOrderConfirmChannelRx:
			if shouldThisMessageBeProcessed(node.receivedMessages, confirmation.SenderID, confirmation.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, confirmation.SenderID, confirmation.MessageID)
				fmt.Printf("%#v \n", confirmation)
			}

		case complete := <-node.orderCompleteChannelRx:
			if shouldThisMessageBeProcessed(node.receivedMessages, complete.SenderID, complete.MessageID) {
				addMessageIDToReceivedMessageMap(node.receivedMessages, complete.SenderID, complete.MessageID)
				fmt.Printf("%#v \n", complete)
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
