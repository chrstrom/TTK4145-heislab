Local order delegation
================
The purpose of this module is to provide an easy way to split incoming local orders into cab and hall orders,
and route them to the channels they are meant for.

### Interface
LocalOrder struct  
Delegator struct  
OrderDelegator(buttonCh <-chan ButtonEvent, cabOrderCh chan<- int, hallOrderCh chan<- LocalOrder)
* Local incoming orders come through the buttonCh channel
* Local cab orders and hall orders are distributed on the cabOrderCh and hallOrderCh respectively.
