package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hazelrah-qr/stim-test/payment"
)

func main() {
	c := make(chan int64, 10)

	client := payment.NewClient(c)
	start := time.Now()

	go func() {
		var sum int64 = 0
		for partial := range c {
			sum += partial
		}
		print("Total: ", sum)
	}()

	filterDate := "2022-02-04"
	wg := &sync.WaitGroup{}

	for i := 1; i <= 10; i++ {
		page := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.ProcessPaymentPage(filterDate, fmt.Sprint(page))
		}()
	}
	wg.Wait()
	close(c)

	fmt.Printf("Processing all data took %s", time.Since(start))
}
