package network

import (
	"fmt"
	"os"

	"github.com/TTK4145-Students-2021/project-gruppe80/network-Go-modul/bcast"
	"github.com/TTK4145-Students-2021/project-gruppe80/network-Go-modul/localip"
	"github.com/TTK4145-Students-2021/project-gruppe80/network-Go-modul/peers"
)

const DuplicatesOfMessages = 1

type ElevatorNode struct {
	id               string
	messageIDCounter int

	peerUpdateChannel                                            chan peers.PeerUpdate
	peerTxEnable                                                 chan bool
	newRequestChannelTx, newRequestChannelRx                     chan NewRequestMessage
	newRequestReplyChannelTx, newRequestReplyChannelRx           chan NewRequestReplyMessage
	delegateOrderChannelTx, delegateOrderChannelRx               chan DelegateOrderMessage
	delegateOrderConfirmChannelTx, delegateOrderConfirmChannelRx chan DelegateOrderConfirmMessage
}

func (node *ElevatorNode) Init() {
	localIP, err := localip.LocalIP()
	if err != nil {
		localIP = "DISCONNECTED"
	}
	node.id = fmt.Sprintf("%v-%v", localIP, os.Getpid())
	fmt.Printf("Init elevator network node with id:%v\n", node.id)

	node.messageIDCounter = 1

	node.peerUpdateChannel = make(chan peers.PeerUpdate)
	node.peerTxEnable = make(chan bool)
	go peers.Transmitter(15647, node.id, node.peerTxEnable)
	go peers.Receiver(15647, node.peerUpdateChannel)

	node.newRequestChannelTx = make(chan NewRequestMessage)
	node.newRequestChannelRx = make(chan NewRequestMessage)
	node.newRequestReplyChannelTx = make(chan NewRequestReplyMessage)
	node.newRequestReplyChannelRx = make(chan NewRequestReplyMessage)
	node.delegateOrderChannelTx = make(chan DelegateOrderMessage)
	node.delegateOrderChannelRx = make(chan DelegateOrderMessage)
	node.delegateOrderConfirmChannelTx = make(chan DelegateOrderConfirmMessage)
	node.delegateOrderConfirmChannelRx = make(chan DelegateOrderConfirmMessage)

	go bcast.Transmitter(20001, node.newRequestChannelTx)
	go bcast.Receiver(20001, node.newRequestChannelRx)

	go bcast.Transmitter(20002, node.newRequestReplyChannelTx)
	go bcast.Receiver(20002, node.newRequestReplyChannelRx)

	go bcast.Transmitter(20003, node.delegateOrderChannelTx)
	go bcast.Receiver(20003, node.delegateOrderChannelRx)

	go bcast.Transmitter(20004, node.delegateOrderConfirmChannelTx)
	go bcast.Receiver(20004, node.delegateOrderConfirmChannelRx)
}

func (node *ElevatorNode) SendNewRequest(floor, direction int) {
	message := NewRequestMessage{node.id, node.messageIDCounter, floor, direction}
	node.messageIDCounter++

	//fmt.Printf("Sending new request%#v \n", messag)

	for i := 0; i < DuplicatesOfMessages; i++ {
		node.newRequestChannelTx <- message
	}
}

func (node *ElevatorNode) SendNewReqestReply(floor, direction, cost int) {
	message := NewRequestReplyMessage{node.id, node.messageIDCounter, floor, direction, cost}
	node.messageIDCounter++

	//fmt.Printf("Sending cost for network request %#v \n", messag)

	for i := 0; i < DuplicatesOfMessages; i++ {
		node.newRequestReplyChannelTx <- message
	}
}

func (node *ElevatorNode) SendDelegateOrder(reciverID string, floor, direction int) {
	message := DelegateOrderMessage{node.id, node.messageIDCounter, reciverID, floor, direction}
	node.messageIDCounter++

	//fmt.Printf("Sending order delegation %#v \n", messag)

	for i := 0; i < DuplicatesOfMessages; i++ {
		node.delegateOrderChannelTx <- message
	}
}

func (node *ElevatorNode) SendDelegateOrderConfirm(reciverID string, floor, direction int) {
	message := DelegateOrderConfirmMessage{node.id, node.messageIDCounter, reciverID, floor, direction}
	node.messageIDCounter++

	for i := 0; i < DuplicatesOfMessages; i++ {
		node.delegateOrderConfirmChannelTx <- message
	}
}
