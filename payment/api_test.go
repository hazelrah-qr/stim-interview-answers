package payment_test

import (
	"testing"

	"github.com/hazelrah-qr/stim-interview-answers/payment"
)

func TestParse(t *testing.T) {
	p, _ := payment.ParsePayment("2022-02-04,name,12")
	if p.Date != "2022-02-04" || p.Name != "name" || p.Amount != 12 {
		t.Errorf("Invalid payment: %v", p)
	}
}

func TestInvalidParse(t *testing.T) {
	_, err := payment.ParsePayment("2022-02-04,name")

	if err == nil {
		t.Errorf("Expected error")
	}
}
