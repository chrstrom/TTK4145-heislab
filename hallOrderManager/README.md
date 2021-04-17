Hall order manager
================
The purpose of this module is to handle hall orders, and is completely separated from cab orders. This is done
since cab orders are exclusively handled locally, while hall orders are shared with every elevator on the network.

The driver part of this module is found in the for-select in OrderManager(). The functionality of this can be
split into three sections:

1. **Local orders**  
If a local order is received, the local elevator is itself responsible for delegating this order to the elevator on the network
with the lowest cost. The elevator will request costs from the other elevators on the network for a set time, and delegate to the
elevator with the lowest cost upon timeout.

2. **Network orders**  
An order from the network can have multiple states. If the order is in the process of being delegated, the elevator will reply to
the sender with the cost of taking this order. If the incoming order has been delegated to the elevator, it will instead take it
and pass it on to the elevator FSM.

3. **Synchronization**  
Since the hall orders are shared on the network, every elevator needs to agree upon who does what.
This is done through an order synchronization system: Upon every order update, the order is shared
with every elevator on the network on a dedicated sync-channel. This synchronization also makes sure
that panel lights are matching for every elevator.


### Interface
OrderManager(ID, <-chan localOrderDelegation.LocalOrder, msg.FSMChannels, msg.NetworkChannels)
* ID is taken as input so that this module and the network module can agree on "who am I?"
* Needs to be run as a goroutine, since it contains a for-select
* Local orders come from the LocalOrder channel
* FSMChannels need to be shared between this call and the elevator FSM routine
* NetworkChannels need to be shared between this call and the network routine
