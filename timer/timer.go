package timer

import "time"

func SendWithDelay(delay time.Duration, ch chan<- int, message int) {
	go sendWithDelayFunction(delay, ch, message)
}

func sendWithDelayFunction(delay time.Duration, ch chan<- int, message int) {
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
