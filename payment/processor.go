package payment

import (
	"strconv"
	"strings"
)

type Processor struct {
	Sum        int64
	FilterDate string
}

type Payment struct {
	Date   string
	Name   string
	Amount int64
}

func NewProcessor(filterDate string) *Processor {
	return &Processor{Sum: 0, FilterDate: filterDate}
}

func ParsePayment(line string) Payment {
	fields := strings.Split(line, ",")
	amount, err := strconv.ParseInt(fields[2], 10, 64)

	if err != nil {
		panic(err)
	}

	return Payment{Date: fields[0], Name: fields[1], Amount: amount}
}

func (proc *Processor) ProcessPayment(p Payment) {
	if p.Date == proc.FilterDate {
		proc.Sum += p.Amount
	}
}
