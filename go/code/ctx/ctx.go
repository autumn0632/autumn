package ctx

import (
	"context"
	"log"
	"time"
)

func DoCancel() {
	context, cancel := context.WithCancel(context.Background())

	go PrintTime(context)

	time.Sleep(10 * time.Second)

	cancel()
	log.Printf("%v", context.Err())
}

func PrintTime(c context.Context) {
	for {
		select {
		case <-c.Done():
			println("1111")
			println( c.Err())
		default:
			time.Sleep(time.Second)
			println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}