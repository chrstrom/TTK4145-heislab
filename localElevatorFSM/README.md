Local elevator FSM
================
The local elevator FSM is responsible for the behavior of each local elevator. Logically, it is separate
from the hall order manager and local order delegator. It simply takes in hall orders and cab orders, and
executes them. It does, however, interface with the hall order manager in that it takes hall orders from it,
and replies with costs when requested. Additionally, it uses the cab order storage module to ensure that no local cab orders are lost.

### Interface

CreateFSMChannelStruct() types.FSMChannels
* Used to create the channels that links the elevator FSM to the hall order manager.

RunElevatorFSM(cabOrderChannel, FSMChannels, NetworkChannels, ioChannels)
* Contains a for-select and should thus be ran as a goroutine
* cabOrderChannel is the channel that links the local delegator to the FSM
* FSMChannels and NetworkChannels link the FSM to the respect
