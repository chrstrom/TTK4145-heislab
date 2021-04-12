package hallOrderManager

import (
	"log"
	"time"

	"../localOrderDelegation"
	msg "../orderTypes"
)

const orderReplyTime = time.Millisecond * 50
const orderDelegationTime = time.Millisecond * 500

type HallOrderManager struct {
	id string

	orders         OrderMap
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork           chan<- msg.OrderStamped
	delegationConfirmToNetwork chan<- msg.OrderStamped
	delegateToNetwork          chan<- msg.OrderStamped
	orderSyncToNetwork         chan<- msg.HallOrder

	replyToRequestFromNetwork         <-chan msg.OrderStamped
	orderDelegationConfirmFromNetwork <-chan msg.OrderStamped
	delegationFromNetwork             <-chan msg.OrderStamped
	orderSyncFromNetwork              <-chan msg.HallOrder

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int

	logger *log.Logger
}
