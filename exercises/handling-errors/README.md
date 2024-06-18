## Exercise #1: Handling Errors

During this exercise, you will:

- Throw and handle exceptions in Temporal Workflows and Activities
- Use non-retryable errors to fail an Activity
- Locate the details of a failure in Temporal Workflows and Activities in the Event History

Make your changes to the code in the `practice` subdirectory (look for
`TODO` comments that will guide you to where you should make changes to
the code). If you need a hint or want to verify your changes, look at
the complete version in the `solution` subdirectory.

## Part A: Throw a non-retryable Application Error to fail an Activity

In this part of the exercise, you will throw a non-retryable Application Failure 
that will fail your Activities.

Application Failures are used to communicate application-specific failures in
Workflows and Activities. In Activities, returning an `ApplicationError` will
cause the Activity to fail. However, unless this Activity is specified as
non-retryable, it will retry according to the Retry Policy. To have an Activity
fail when an `ApplicationError` is returned, set it as non-retryable. Any other
error that is returned in Go is automatically converted to an `ActivityError`.

1. Start by opening and reviewing the `activities.go` file, familiarizing yourself with
   the Activity functions.
2. In the `SendBill` Activity, notice how an error is thrown if the
   `chargeAmount` is less than 0. If the calculated amount to charge the
   customer is negative, the Activity returns with a non-`nil` error.
   Specifically, it returns with `fmt.Errorf("invalid charge amount: %d (< 1)",
   chargeAmount)`. In Go, returning `fmt.Errorf()` lets you return an error with
   customized output without having to define a new custom error class — it will
   be converted to an `ActivityError` by Temporal. By default, depending on
   this Workflow's retry policy, this Activity may be retried if it fails, but
   since an invalid charge amount would probably reflect bad data being
   provided, you may not want it to retry. To do that, you can replace the
   `fmt.Errorf()` return value with
   `temporal.NewNonRetryableApplicationError(fmt.Sprintf("invalid charge amount: %d (< 1)", chargeAmount), "invalidChargeError", nil, nil)`. This will use the same formatting, but will
   mark the error as a non-retryable `invalidChargeError`.
3. Go to `ProcessCreditCard` Activity, where you will return an
   `ApplicationError` if the credit card fails its validation step. In this
   Activity, you will throw an error if the entered credit card number does not
   have 16 digits. Use the `NewNonRetryableApplicationError()` code from the
   previous step as a reference. Designate this as a `CreditCardError`.
4. Save your file.

## Part B: Catch the Activity Failure

In this part of the exercise, you will catch the `ApplicationError` that was
thrown from the `ProcessCreditCard` Activity and handle it.

1. Open `workflow.go` in your text editor.
2. `ProcessCreditCard` is run like so:
   `err = workflow.ExecuteActivity(ctx, ProcessCreditCard, confirmation).Get(ctx, &chargestatus)`.
   Note that the results are returned to `err`, and the next code block, `if err != nil {}`,
   handles any errors that have been returned from the Activity. By default, if a non-retryable
   `ApplicationFailure` is thrown, the Workflow will fail. However, it is possible to catch this
   failure and either handle it or continue to propagate it up.
3. Within the `if err != nil {}` block, check for an `ApplicationError` like so:
   ```
   var applicationErr *temporal.ApplicationError
		if errors.As(err, &applicationErr) {}
   ```			
   Within that `if errors.As() {}` block, you can handle the error, add additional logging, and so on. For example, you could add:
   ```
   println("Billing timestamp of failed order:", confirmation.BillingTimestamp)
   logger.Error("Unable to charge credit card", "Error", err)
   ```
   Finally, make sure that you add `return OrderConfirmation{}, err` at the end of the
   `if err != nil {}` block, so the error is returned.
4. Save your file.

## Part C: Run the Workflow

In this part of the exercise, you will run your Workflow and see both your
Workflow and Activity succeed and fail.

The starter used to run this Workflow, `start/main.go`, contains a "valid"
16-digit credit card number which will be used by default, causing the Workflow
to complete successfully.

**First, run the Workflow successfully:**

1. In one terminal, start the Worker by running:
   ```bash
   go run worker/main.go
   ```
2. In another terminal, start the Workflow by executing `start/main.go`:
   ```bash
   go run start/main.go
   ```
3. In the Web UI, verify that the Workflow ran successfully to completion.

**Next, run the Workflow with a bad credit card number to trigger a failure:**

1. You do not need to restart the Worker. In the terminal where you ran
   `start/main.go`, run it again with an additional parameter, `--creditcard`.
   Supply an obviously wrong credit card number like "1234":
   ```bash
   go run start/main.go --creditcard 1234
   ```
2. You should see the Workflow fail in the terminal where you executed `go run start/main.go`.
   You can also check the Web UI and view the failure there.

### This is the end of the exercise.