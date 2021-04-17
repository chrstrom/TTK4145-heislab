Network
================
The purpose of this module is to link elevators together through a P2P connection.
Every elevator on the network can talk to every other elevator through UDP broadcast,
although this module also allows for addressing through elevator node ID.


### Interface
 GetNodeID() string   
 CreateNetworkChannelStruct() msg.NetworkChannels  
 NetworkNode(ID, FSMChannels, NetworkChannels)
 * Contains a for-select, and should thus be called as a goroutine.
