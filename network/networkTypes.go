package network

import (
	"log"

	msg "../messageTypes"
	"../network/peers"
)

type Node struct {
	id               string
	messageIDCounter int

	networkChannels msg.NetworkChannels

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
