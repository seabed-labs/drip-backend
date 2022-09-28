package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func generateNumbers(total int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for idx := 1; idx <= total; idx++ {
		val := idx
		ch <- val
	}
}

func printNumbers(ctx context.Context, ch <-chan int, wg *sync.WaitGroup, name string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg2 := <-ch:
			fmt.Printf("read %d from channel %s \n", msg2, name)
		}
	}
}

func main() {
	c := make(chan int)
	var w sync.WaitGroup
	w.Add(5)

	for i := 1; i <= 5; i++ {
		go func(i int, ci <-chan int) {
			j := 1
			for v := range ci {
				time.Sleep(time.Millisecond)
				fmt.Printf("%d.%d got %d\n", i, j, v)
				j += 1
			}
			w.Done()
		}(i, c)
	}

	for i := 1; i <= 25; i++ {
		c <- i
	}
	close(c)
	w.Wait()
}
