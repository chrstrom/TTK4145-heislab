package network

type NewRequestNetworkMessage struct {
	SenderID         string
	MessageID        int
	Floor, Direction int
	OrderID          int
}

type NewRequestReplyNetworkMessage struct {
	SenderID               string
	MessageID              int
	ReceiverID             string
	Floor, Direction, Cost int
	OrderID                int
}

type DelegateOrderNetworkMessage struct {
	SenderID         string
	MessageID        int
	ReceiverID       string
	Floor, Direction int
	OrderID          int
}

type DelegateOrderConfirmNetworkMessage struct {
	SenderID         string
	MessageID        int
	ReceiverID       string
	Floor, Direction int
	OrderID          int
}

type OrderCompleteNetworkMessage struct {
	SenderID         string
	MessageID        int
	ReceiverID       string
	Floor, Direction int
	OrderID          int
}
