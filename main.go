package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/hazelrah-qr/stim-interview-answers/payment"
	"golang.org/x/sync/semaphore"
)

const url string = "https://stim-codetest-stimseassets.ew.r.appspot.com/"

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:8081", nil))
	}()

	c := make(chan int64, 10)

	client := payment.NewClient(url, c)
	now := time.Now()

	go func() {
		var sum int64 = 0
		for partial := range c {
			sum += partial
		}
		fmt.Printf("Total: %d\n", sum)
	}()

	filter := "2022-02-14"
	wg := &sync.WaitGroup{}
	sem := semaphore.NewWeighted(2)

	// TODO: Pass more contexts
	for i := 1; i <= 100; i++ {
		page := i
		wg.Add(1)
		sem.Acquire(context.Background(), 1)
		go func() {
			defer wg.Done()
			defer sem.Release(1)
			client.ProcessPaymentPage(filter, fmt.Sprint(page))
		}()
	}

	wg.Wait()
	close(c)

	// Make deterministic - wait for goroutine summing all partials
	fmt.Printf("Processing all data took %s", time.Since(now))
}
