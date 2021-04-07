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
	ID            int
	State         OrderStateType
	Floor, Dir    int
	costs         map[string]int
	OwnerID       string
	DelegatedToID string
}

type HallOrderManager struct {
	id string

	orders         OrderMap
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork  	chan<- network.OrderStamped
	delegateToNetwork 	chan<- network.OrderStamped
	orderSyncToNetwork 	chan<- network.OrderSync

	requestReplyFromNetwork           	<-chan network.OrderStamped
	orderDelegationConfirmFromNetwork 	<-chan network.OrderStamped
	orderSyncFromNetwork				<-chan network.OrderSync

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
}
