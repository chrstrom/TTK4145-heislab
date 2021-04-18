package timer

import (
	"time"

	msg "../orderTypes"
)

func SendWithDelayInt(delay time.Duration, ch chan<- int, message int) {
	go sendWithDelayIntFunction(delay, ch, message)
}

func sendWithDelayIntFunction(delay time.Duration, ch chan<- int, message int) {
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
