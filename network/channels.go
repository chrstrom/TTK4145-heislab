package network

type NetworkChannels struct {
	RequestToNetwork           chan NewRequest
	DelegateOrderToNetwork     chan Delegation
	RequestReplyToNetwork      chan RequestReply
	DelegationConfirmToNetwork chan DelegationConfirm
	OrderCompleteToNetwork     chan OrderComplete

	RequestFromNetwork           chan NewRequest
	DelegateFromNetwork          chan Delegation
	RequestReplyFromNetwork      chan RequestReply
	DelegationConfirmFromNetwork chan DelegationConfirm
	OrderCompleteFromNetwork     chan OrderComplete
}
