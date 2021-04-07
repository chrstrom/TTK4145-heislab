package network

type NetworkChannels struct {
	RequestToNetwork           chan Order
	DelegateOrderToNetwork     chan Order
	RequestReplyToNetwork      chan Order
	DelegationConfirmToNetwork chan Order
	OrderCompleteToNetwork     chan Order
	SyncOrderToNetwork		   chan OrderSync

	RequestFromNetwork           chan Order
	DelegateFromNetwork          chan Order
	RequestReplyFromNetwork      chan Order
	DelegationConfirmFromNetwork chan Order
	OrderCompleteFromNetwork     chan Order
	SyncOrderFromNetwork		 chan OrderSync
}


func CreateNetworkChannelStruct() NetworkChannels {
	var networkChannels NetworkChannels

	networkChannels.RequestToNetwork 				= make(chan Order)
	networkChannels.DelegateOrderToNetwork 			= make(chan Order)
	networkChannels.RequestReplyToNetwork 			= make(chan Order)
	networkChannels.DelegationConfirmToNetwork 		= make(chan Order)
	networkChannels.OrderCompleteToNetwork 			= make(chan Order)
	networkChannels.SyncOrderToNetwork				= make(chan OrderSync)

	networkChannels.RequestFromNetwork 				= make(chan Order)
	networkChannels.DelegateFromNetwork 			= make(chan	Order)
	networkChannels.RequestReplyFromNetwork 		= make(chan Order)
	networkChannels.DelegationConfirmFromNetwork 	= make(chan Order)
	networkChannels.OrderCompleteFromNetwork 		= make(chan Order)
	networkChannels.SyncOrderFromNetwork			= make(chan OrderSync)

	return networkChannels
}