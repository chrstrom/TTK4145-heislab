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


func CreateNetworkChannelStruct() NetworkChannels {
	var networkChannels NetworkChannels

	networkChannels.RequestToNetwork 				= make(chan NewRequest)
	networkChannels.DelegateOrderToNetwork 			= make(chan Delegation)
	networkChannels.RequestReplyToNetwork 			= make(chan RequestReply)
	networkChannels.DelegationConfirmToNetwork 		= make(chan DelegationConfirm)
	networkChannels.OrderCompleteToNetwork 			= make(chan OrderComplete)
	networkChannels.RequestFromNetwork 				= make(chan NewRequest)
	networkChannels.DelegateFromNetwork 			= make(chan	Delegation)
	networkChannels.RequestReplyFromNetwork 		= make(chan RequestReply)
	networkChannels.DelegationConfirmFromNetwork 	= make(chan DelegationConfirm)
	networkChannels.OrderCompleteFromNetwork 		= make(chan OrderComplete)

	return networkChannels
}