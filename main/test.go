package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	defer func() {
		fmt.Printf("defer running ...\n")
	}()

	go func() {
		fmt.Printf("test ...\n")
		wg.Done()
		wg.Done()
		return

	}()

	wg.Wait()
}
