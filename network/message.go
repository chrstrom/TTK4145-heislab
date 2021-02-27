package network

type NewRequest struct {
	ID         string
	OrderID    int
	Floor, Dir int
}

type RequestReply struct {
	ID         string
	OrderID    int
	Floor, Dir int
	Cost       int
}

type Delegation struct {
	ID         string
	OrderID    int
	Floor, Dir int
}

type DelegationConfirm struct {
	ID         string
	OrderID    int
	Floor, Dir int
}

type OrderComplete struct {
	ID         string
	OrderID    int
	Floor, Dir int
}
