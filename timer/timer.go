package timer

import "time"

func SendWithDelay(delay time.Duration, ch chan<- int, message int) {
	go sendWithDelayFunction(delay, ch, message)
}

func sendWithDelayFunction(delay time.Duration, ch chan<- int, message int) {
	<-time.After(delay)
	ch <- message
}
