package messageTypes

import "../network/peers"

type NetworkOrder struct {
	SenderID   string
	MessageID  int
	ReceiverID string
	Order      OrderStamped
}

type NetworkHallOrder struct {
	SenderID   string
	MessageID  int
	ReceiverID string
	Order      HallOrder
}

type NetworkChannels struct {
	RequestToNetwork           chan OrderStamped
	DelegateOrderToNetwork     chan OrderStamped
	DelegationConfirmToNetwork chan OrderStamped
	OrderCompleteToNetwork     chan OrderStamped
	SyncOrderToNetwork         chan HallOrder

	DelegateFromNetwork          chan OrderStamped
	ReplyToRequestFromNetwork    chan OrderStamped
	DelegationConfirmFromNetwork chan OrderStamped
	OrderCompleteFromNetwork     chan OrderStamped
	SyncOrderFromNetwork         chan HallOrder
	PeerUpdate                   chan peers.PeerUpdate
}
