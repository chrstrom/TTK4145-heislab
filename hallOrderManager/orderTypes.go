package hallOrderManager

import (
	"log"

	"../elevio"
	"../localOrderDelegation"
	msg "../messageTypes"
	"../network/peers"
)

type HallOrderManager struct {
	id string

	orders         OrderMap
	orderIDCounter int

	localRequestChannel <-chan localOrderDelegation.LocalOrder

	requestToNetwork           chan<- msg.OrderStamped
	delegationConfirmToNetwork chan<- msg.OrderStamped
	delegateToNetwork          chan<- msg.OrderStamped
	orderSyncToNetwork         chan<- msg.HallOrder
	delegateToLocalElevator    chan<- elevio.ButtonEvent
	requestElevatorCost        chan<- msg.RequestCost

	replyToRequestFromNetwork         <-chan msg.OrderStamped
	orderDelegationConfirmFromNetwork <-chan msg.OrderStamped
	delegationFromNetwork             <-chan msg.OrderStamped
	orderSyncFromNetwork              <-chan msg.HallOrder
	peerUpdateChannel                 <-chan peers.PeerUpdate
	elevatorCost                      <-chan int
	orderComplete                     <-chan elevio.ButtonEvent
	ElevatorUnavailable               <-chan bool

	orderReplyTimeoutChannel      chan int
	orderDelegationTimeoutChannel chan int
	orderCompleteTimeoutChannel   chan msg.HallOrder

	logger *log.Logger
}
