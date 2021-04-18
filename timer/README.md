Timer
================
The timer module in place since we need capabiltiies beyond the builtin golang timer, namely
to send a message with type that we have defined ourselves after a specified amount of time
has elapsed.

### Interface
SendWithDelayInt(time.Duration, chan<- int, int)  
SendWithDelayHallOrder(time.Duration, chan<- HallOrder, HallOrder)
