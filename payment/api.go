package payment

import (
	"bufio"
	"net/http"
	"os"
	"time"
)

const url string = "https://stim-codetest-stimseassets.ew.r.appspot.com/"

type StimPaymentClient struct {
	client     *http.Client
	processing chan<- int64
}

func NewClient(c chan<- int64) *StimPaymentClient {
	client := &http.Client{Timeout: time.Second * 4}
	return &StimPaymentClient{client: client, processing: c}
}

func (s *StimPaymentClient) ProcessPaymentPage(filterDate string, page string) {
	req, err := http.NewRequest(http.MethodGet, url+page, nil)

	if err != nil {
		panic(err)
	}

	resp, err := s.client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(err)
	}

	processor := NewProcessor(filterDate)

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || text == "date,name,amount" {
			continue
		}
		processor.ProcessPayment(ParsePayment(text))
	}
	s.processing <- processor.Sum
}

func (s *StimPaymentClient) ProcessPaymentFile(filterDate string, page string) {
	file, _ := os.Open("data/" + page + ".csv")

	defer file.Close()
	scanner := bufio.NewScanner(file)
	processor := NewProcessor(filterDate)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || text == "date,name,amount" {
			continue
		}
		processor.ProcessPayment(ParsePayment(text))
	}
	s.processing <- processor.Sum
}
