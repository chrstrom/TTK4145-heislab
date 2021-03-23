package hallOrderManager

import (
	"time"

	"../localOrderDelegation"
	"../network"
)

type OrderStateType int

const (
	Received OrderStateType = iota
	Delegate
	Serving
)
const orderReplyTime = time.Millisecond * 50
const orderDelegationTime = time.Millisecond * 50

type Order struct {
	State         OrderStateType
	Floor, Dir    int
	costs         map[string]int
	DelegatedToID string
}

type HallOrderManager struct {
	id string

	orders         map[int]Order
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork  chan<- network.NewRequest
	delegateToNetwork chan<- network.Delegation

	requestReplyFromNetwork           <-chan network.RequestReply
	orderDelegationConfirmFromNetwork <-chan network.DelegationConfirm

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
}
