package network


type Order struct {
	Floor	int
	Dir		int
	Cost    int
}

type OrderStamped struct {
	ID string
	OrderID int
	Order Order
}


type NetworkOrder struct {
	SenderID               string
	MessageID              int
	ReceiverID             string
	Order		OrderStamped
}

type OrderSync struct {
	ID         		string
	OrderID    		int
	Floor, Dir		int
	DelegatedToID 	int
}


type OrderSyncNetworkMessage struct {
	SenderID         string
	MessageID        int
	ReceiverID       string
	Order 			 OrderSync
}

type NetworkChannels struct {
	RequestToNetwork           chan OrderStamped
	DelegateOrderToNetwork     chan OrderStamped
	RequestReplyToNetwork      chan OrderStamped
	DelegationConfirmToNetwork chan OrderStamped
	OrderCompleteToNetwork     chan OrderStamped
	SyncOrderToNetwork		   chan OrderSync

	RequestFromNetwork           chan OrderStamped
	DelegateFromNetwork          chan OrderStamped
	RequestReplyFromNetwork      chan OrderStamped
	DelegationConfirmFromNetwork chan OrderStamped
	OrderCompleteFromNetwork     chan OrderStamped
	SyncOrderFromNetwork		 chan OrderSync
}
