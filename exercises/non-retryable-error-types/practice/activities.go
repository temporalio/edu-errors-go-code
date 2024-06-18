package pizza

import (
	"context"
	"fmt"
	"math/rand"
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
			// TODO Part A: Change this to a `NewApplicationError` so it's retryable.
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

func NotifyDeliveryDriver(ctx context.Context) error {
	/* This is a simulation of attempting to notify a delivery driver that
	the order is ready for delivery. It starts by generating a number from 0 - 14.
	From there a loop is iterated over from 0 < 10, each time checking to
	see if the random number matches the loop counter and then sleeping for 5
	seconds. Each iteration of the loop sends a heartbeat back letting the
	Workflow know that progress is still being made. If the number matches a
	loop counter, it is a success. If it doesn't, then a delivery driver was
	unable to be contacted and failure is returned.
	*/

	logger := activity.GetLogger(ctx)
	SuccessSimulation := rand.Intn(15)

	for x := 0; x < 10; x++ {
		if SuccessSimulation == x {
			logger.Info("Delivery driver responded")
			return nil
		}
		// TODO Part C: Add a call to `activity.RecordHeartbeat()`
		logger.Info("Heartbeat:", x)
		// TODO Part F: Lengthen the `time.Sleep()` call so the Activity fails.
		time.Sleep(time.Second * 5)
	}

	return temporal.NewApplicationError(fmt.Sprintf("Driver didn't respond."), "DriverError")
}
