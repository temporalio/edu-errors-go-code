package pizza

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

func GetDistance(ctx context.Context, address Address) (Distance, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetDistance invoked; determining distance to customer address")

	// this is a simulation, which calculates a fake (but consistent)
	// distance for a customer address based on its length. The value
	// will therefore be different when called with different addresses,
	// but will be the same across all invocations with the same address.
	kilometers := len(address.Line1) + len(address.Line2) - 10
	if kilometers < 1 {
		kilometers = 5
	}

	distance := Distance{
		Kilometers: kilometers,
	}

	logger.Debug("GetDistance complete", "Distance", distance.Kilometers)
	return distance, nil
}

func SendBill(ctx context.Context, bill Bill) (OrderConfirmation, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SendBill invoked", "Customer", bill.CustomerID, "Amount", bill.Amount)

	chargeAmount := bill.Amount

	// This month's special offer: Get $5 off all orders over $30
	if bill.Amount > 3000 {
		logger.Info("Applying discount")

		chargeAmount -= 500 // reduce amount charged by 500 cents
	}

	// reject invalid amounts before calling the payment processor
	if chargeAmount < 0 {
		return OrderConfirmation{},
			temporal.NewNonRetryableApplicationError(fmt.Sprintf("invalid charge amount: %d (< 1)", chargeAmount), "invalidChargeError", nil, nil)
	}

	// pretend we called a payment processing service here :-)

	confirmation := OrderConfirmation{
		OrderNumber:        bill.OrderNumber,
		ConfirmationNumber: "AB9923",
		Status:             "SUCCESS",
		BillingTimestamp:   time.Now().Unix(),
		Amount:             chargeAmount,
	}

	logger.Debug("SendBill complete", "ConfirmationNumber", confirmation.ConfirmationNumber)

	return confirmation, nil
}

func ProcessCreditCard(ctx context.Context, address Address) (ChargeStatus, error) {
	// pretend to charge card here
	chargestatus := ChargeStatus{
		Success: true,
	}

	if len(address.CardNumber) != 16 {
		return chargestatus, temporal.NewNonRetryableApplicationError("Credit Card Charge Error", "CreditCardError", nil, nil)
	} else {
		return chargestatus, nil
	}
}

func UpdateInventory(ctx context.Context, items []Pizza) error {
	// Here you would call your inventory management system to reduce the stock of your pizza inventory
	logger := activity.GetLogger(ctx)
	logger.Info("UpdateInventory complete", items)
	return nil
}

func RevertInventory(ctx context.Context, items []Pizza) error {
	// Here you would call your inventory management system to add the ingredients back into the pizza inventory.
	logger := activity.GetLogger(ctx)
	logger.Info("RevertInventory complete", items)
	return nil
}

func RefundCustomer(ctx context.Context, bill Bill) error {
	// Simulate refunding the customer
	logger := activity.GetLogger(ctx)
	logger.Info("Refunding", bill.Amount, "to customer", bill.CustomerID, "for order", bill.OrderNumber)
	return nil
}
