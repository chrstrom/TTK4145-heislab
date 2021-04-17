Network
================
The purpose of this module is to link elevators together through a P2P connection.
Every elevator on the network can talk to every other elevator through UDP broadcast,
although this module also allows for addressing through elevator node ID.

The network node can be broken down into two categories:

1. **Channels connecting the elevators together on the network**  
Whenever a message is to be sent out from the local elevator to the network, it 
will pass through one of these channels. Here, a network version of the outgoing 
message is generated. Multiple copies as defined by the config file are then sent out
to the network using the handout UDP broadcast package.


2. **Channels connecting the hall order manager to this module**  
The main purpose of this part is to pass on messages from other elevators to the
local elevator. Some basic filtering is done to avoid handling invalid or duplicate
messages. Before a message is sent to the hall order manager, it is added to a map 
of received messages, to aid in the filtering process.


### Interface
 GetNodeID() string   
 CreateNetworkChannelStruct() msg.NetworkChannels  
 NetworkNode(ID, FSMChannels, NetworkChannels)
 * Contains a for-select, and should thus be called as a goroutine.
