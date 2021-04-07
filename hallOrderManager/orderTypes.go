package hallOrderManager

import (
	"time"

	"../localOrderDelegation"
	"../network"
)

const orderReplyTime = time.Millisecond * 50
const orderDelegationTime = time.Millisecond * 50

type HallOrderManager struct {
	id string

	orders         OrderMap
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork           chan<- network.OrderStamped
	delegationConfirmToNetwork chan<- network.OrderStamped
	delegateToNetwork          chan<- network.OrderStamped
	orderSyncToNetwork         chan<- network.HallOrder

	requestReplyFromNetwork           <-chan network.OrderStamped
	orderDelegationConfirmFromNetwork <-chan network.OrderStamped
	delegationFromNetwork             <-chan network.OrderStamped
	orderSyncFromNetwork              <-chan network.HallOrder

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
}
