package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	groupWait sync.WaitGroup
)

func main() {
	groupWait.Add(1)
	go t1()
	groupWait.Add(1)
	go t2()
	groupWait.Wait()
}

func t1() {
	fmt.Println("start")
	time.Sleep(3 * time.Second)
	groupWait.Done()
	fmt.Println("end")
}

func t2() {
	fmt.Println("start")
	time.Sleep(1 * time.Second)
	groupWait.Done()
	fmt.Println("end")
}
