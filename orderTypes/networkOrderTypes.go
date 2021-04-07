package orderTypes

type NetworkOrder struct {
	SenderID   string
	MessageID  int
	ReceiverID string
	Order      OrderStamped
}

type NetworkChannels struct {
	RequestToNetwork           chan OrderStamped
	DelegateOrderToNetwork     chan OrderStamped
	RequestReplyToNetwork      chan OrderStamped
	DelegationConfirmToNetwork chan OrderStamped
	OrderCompleteToNetwork     chan OrderStamped
	SyncOrderToNetwork         chan HallOrder

	RequestFromNetwork           chan OrderStamped
	DelegateFromNetwork          chan OrderStamped
	RequestReplyFromNetwork      chan OrderStamped
	DelegationConfirmFromNetwork chan OrderStamped
	OrderCompleteFromNetwork     chan OrderStamped
	SyncOrderFromNetwork         chan HallOrder
}
