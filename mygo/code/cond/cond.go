package cond

import (
	"sync"
	"time"
)

func CondPrintf() {
	c := sync.NewCond(&sync.Mutex{})
	wg := sync.WaitGroup{}
	wg.Add(11)

	for i := 0; i < 10; i++ {
		go NumPrint(&wg, c, i)
	}

	time.Sleep(time.Second)
	go func() {
		defer wg.Done()
		println("i must run first,then print num")
		CondBroadCast(c)
	}()

	wg.Wait()
}

func NumPrint(wg *sync.WaitGroup, c *sync.Cond, i int) {
	defer wg.Done()
	c.L.Lock()
	c.Wait()
	println(i)
	c.L.Unlock()
}

func CondSignal(c *sync.Cond) {
	c.L.Lock()
	c.Signal()
	c.L.Unlock()
}

func CondBroadCast(c *sync.Cond) {
	c.L.Lock()
	c.Broadcast()
	c.L.Unlock()
}
