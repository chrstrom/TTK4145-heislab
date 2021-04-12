package hallOrderManager

import (
	"time"

	"../elevio"
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
	OwnerID       string
	ID            int
	DelegatedToID string
	State         OrderStateType
	Floor, Dir    int
	costs         map[string]int
}

type HallOrderManager struct {
	id string

	orders         OrderMap
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork           chan<- network.NewRequest
	delegateToNetwork          chan<- network.Delegation
	delegationConfirmToNetwork chan<- network.DelegationConfirm
	delegateToLocalElevator    chan<- elevio.ButtonEvent
	requestElevatorCost        chan<- elevio.ButtonEvent

	requestReplyFromNetwork           <-chan network.RequestReply
	orderDelegationConfirmFromNetwork <-chan network.DelegationConfirm
	delegationFromNetwork             <-chan network.Delegation
	elevatorCost                      <-chan int
	orderComplete                     <-chan elevio.ButtonEvent

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
}
