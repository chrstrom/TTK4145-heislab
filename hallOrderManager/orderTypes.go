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

type HallOrder struct {
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

	requestToNetwork  			chan<- network.OrderStamped
	delegationConfirmToNetwork 	chan<- network.OrderStamped
	delegateToNetwork 			chan<- network.OrderStamped
	orderSyncToNetwork 			chan<- 	HallOrder

	requestReplyFromNetwork           	<-chan network.OrderStamped
	orderDelegationConfirmFromNetwork	<-chan network.OrderStamped
	delegationFromNetwork             	<-chan network.OrderStamped
	orderSyncFromNetwork  				<-chan HallOrder

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
}
