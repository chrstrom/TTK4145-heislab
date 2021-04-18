Message types
================
This package was introduced to avoid circular dependencies in the other packages,
and contains types for the following modules:

### FSM
* HallOrderManager const
* Network const
* RequestCost struct
* FSMChannels struct

### Hall order manager
* OrderStateType int
* Order struct
* OrderStamped struct
* HallOrder struct

### Network
* NetworkOrder
* NetworkHallOrder
* NetworkChannels struct
