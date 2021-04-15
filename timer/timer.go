package timer

import (
	"time"

	msg "../orderTypes"
)

func SendWithDelay(delay time.Duration, ch chan<- int, message int) {
	go sendWithDelayFunction(delay, ch, message)
}

func sendWithDelayFunction(delay time.Duration, ch chan<- int, message int) {
	<-time.After(delay)
	ch <- message
}

func SendWithDelayHallOrder(delay time.Duration, ch chan<- msg.HallOrder, message msg.HallOrder) {
	go sendWithDelayHallOrderFunction(delay, ch, message)
}

func sendWithDelayHallOrderFunction(delay time.Duration, ch chan<- msg.HallOrder, message msg.HallOrder) {
	<-time.After(delay)
	ch <- message
}

func FsmSendWithDelay(delay time.Duration, ch chan<- int) {
	go fsmSendWithDelayFunction(delay, ch)
}

func fsmSendWithDelayFunction(delay time.Duration, ch chan<- int) {
	ch <- 1
	<-time.After(delay)
	ch <- -1
}
