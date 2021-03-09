package timer

import "time"

func SendWithDelay(delay time.Duration, ch chan<- int, message int) {
	go sendWithDelayFunction(delay, ch, message)
}

func sendWithDelayFunction(delay time.Duration, ch chan<- int, message int) {
	<-time.After(delay)
	ch <- message
}

func FsmSendWithDelay(delay time.Duration, ch chan<- int, reset *int) {
	*reset++
	go fsmSendWithDelayFunction(delay, ch, reset)
}

func fsmSendWithDelayFunction(delay time.Duration, ch chan<- int, reset *int) {
	<-time.After(delay)
	*reset--
	if *reset == 0 {
		ch <- 1
	}
}
