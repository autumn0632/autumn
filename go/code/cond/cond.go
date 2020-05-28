package cond

import (
	"sync"
	"time"
)

func CondPrintf() {
	//c := sync.NewCond(&sync.Mutex{})
	wg := sync.WaitGroup{}
	wg.Add(11)

	for i := 0; i < 10; i++ {
		//go NumPrint(&wg, c, i)
		go NumPrint2(&wg, i)
	}

	time.Sleep(time.Second)
	go func() {
		defer wg.Done()
		println("i must run first,then print num")
		//CondBroad(c)
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

func NumPrint2(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	println(i)
}

func CondBroad(c *sync.Cond) {
	c.L.Lock()
	c.Broadcast()
	c.L.Unlock()
}
