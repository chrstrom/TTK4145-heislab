package network


type Order struct {
	ID         string
	OrderID    int
	Floor, Dir int
	Cost       int
}


type NetworkOrder struct {
	SenderID               string
	MessageID              int
	ReceiverID             string
	Floor, Direction, Cost int
	OrderID                int
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