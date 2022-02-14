package payment

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type StimApiClient struct {
	client *http.Client
	url    string
	sum    chan<- int64
}

// Create a new stim api client
func NewClient(url string, c chan<- int64) *StimApiClient {
	client := &http.Client{Timeout: time.Second * 6}
	return &StimApiClient{client: client, sum: c, url: url}
}

// Process single page of payments
func (s *StimApiClient) ProcessPaymentPage(filterDate string, page string) {
	req, err := http.NewRequest(http.MethodGet, s.url+page, nil)

	if err != nil {
		panic(err)
	}

	resp, err := s.client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to execute request %s (%s)", s.url+page, resp.Status)
		os.Exit(1)
	}

	s.process(resp.Body, filterDate)
}

// Process single file of payments
func (s *StimApiClient) ProcessPaymentFile(filterDate string, page string) {
	file, _ := os.Open("data/" + page + ".csv")
	defer file.Close()
	s.process(file, filterDate)
}

func (s *StimApiClient) process(body io.Reader, filter string) {
	scanner := bufio.NewScanner(body)
	processor := NewProcessor(filter)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || text == "date,name,amount" {
			continue
		}

		p, err := ParsePayment(text)

		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		processor.ProcessPayment(p)
	}
	s.sum <- processor.Sum
}
