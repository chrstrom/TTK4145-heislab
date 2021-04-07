package network

type NetworkChannels struct {
	RequestToNetwork           chan OrderStamped
	DelegateOrderToNetwork     chan OrderStamped
	RequestReplyToNetwork      chan OrderStamped
	DelegationConfirmToNetwork chan OrderStamped
	OrderCompleteToNetwork     chan OrderStamped
	SyncOrderToNetwork		   chan OrderSync

	RequestFromNetwork           chan OrderStamped
	DelegateFromNetwork          chan OrderStamped
	RequestReplyFromNetwork      chan OrderStamped
	DelegationConfirmFromNetwork chan OrderStamped
	OrderCompleteFromNetwork     chan OrderStamped
	SyncOrderFromNetwork		 chan OrderSync
}


func CreateNetworkChannelStruct() NetworkChannels {
	var networkChannels NetworkChannels

	networkChannels.RequestToNetwork 				= make(chan OrderStamped)
	networkChannels.DelegateOrderToNetwork 			= make(chan OrderStamped)
	networkChannels.RequestReplyToNetwork 			= make(chan OrderStamped)
	networkChannels.DelegationConfirmToNetwork 		= make(chan OrderStamped)
	networkChannels.OrderCompleteToNetwork 			= make(chan OrderStamped)
	networkChannels.SyncOrderToNetwork				= make(chan OrderSync)

	networkChannels.RequestFromNetwork 				= make(chan OrderStamped)
	networkChannels.DelegateFromNetwork 			= make(chan	OrderStamped)
	networkChannels.RequestReplyFromNetwork 		= make(chan OrderStamped)
	networkChannels.DelegationConfirmFromNetwork 	= make(chan OrderStamped)
	networkChannels.OrderCompleteFromNetwork 		= make(chan OrderStamped)
	networkChannels.SyncOrderFromNetwork			= make(chan OrderSync)

	return networkChannels
}