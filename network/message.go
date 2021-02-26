package network

type NewRequestMessage struct {
	SenderID         string
	MessageID        int
	Floor, Direction int
}

type NewRequestReplyMessage struct {
	SenderID               string
	MessageID              int
	Floor, Direction, Cost int
}

type DelegateOrderMessage struct {
	SenderID         string
	MessageID        int
	ReceiverID       string
	Floor, Direction int
}

type DelegateOrderConfirmMessage struct {
	SenderID         string
	MessageID        int
	ReceiverID       string
	Floor, Direction int
}
