package pizza

import (
	"errors"
	"fmt"
)

const TaskQueueName = "pizza-tasks"

type Address struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	CardNumber string
}

type Customer struct {
	CustomerID int
	Name       string
	Email      string
	Phone      string
}

type Pizza struct {
	Description string
	Price       int
}

type PizzaOrder struct {
	OrderNumber string
	Customer    Customer
	Items       []Pizza
	IsDelivery  bool
	Address     Address
}

type Distance struct {
	Kilometers int
}

type ChargeStatus struct {
	Success bool
}

type ChargeError struct {
	StatusCode int
	Err        error
}

func (e *ChargeError) Error() string {
	return fmt.Sprintf("status %d: err %v", e.StatusCode, e.Err)
}

func chargeRequestError() error {
	return &ChargeError{
		StatusCode: 503,
		Err:        errors.New("Credit Card Charge Error"),
	}
}

type Bill struct {
	CustomerID  int
	OrderNumber string
	Description string
	Amount      int
}

type OrderConfirmation struct {
	OrderNumber        string
	Status             string
	ConfirmationNumber string
	BillingTimestamp   int64
	Amount             int
}
