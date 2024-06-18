## Exercise #3: Rollback with the Saga Pattern

During this exercise, you will:

- Orchestrate Activities using a Saga pattern to implement compensating transactions
- Handle failures with rollback logic

Make your changes to the code in the `practice` subdirectory (look for `TODO`
comments that will guide you to where you should make changes to the code). If
you need a hint or want to verify your changes, look at the complete version in
the `solution` subdirectory.

## Part A: Review your new rollback Activities

This Exercise uses the same structure, and the same Error, as in the previous
Exercises — meaning that it will fail at the very end on `ProcessCreditCard` if
you provide it with a bad credit card number.

Three new Activities have been created to demonstrate rollback actions.
`UpdateInventory` is a new step that would run normally (whether or not the
Workflow encounters an error). `RevertInventory` has also been added as a
compensating action for `UpdateInventory`. Finally, `RefundCustomer` has been
added as a compensating action for `SendBill`.

1. Review these new Activities at the end of the `activities.go` file. None of
   them make actual inventory or billing changes, because the intent of this
   Activity is to show Temporal features, but you should be able to see where
   you could add functionality here.
2. Close the file.

## Part B: Add your new rollback Activities to your Workflow

1. Open `workflow.go`. Note that the `UpdateInventory` Activity has been added
   to your Workflow after validating an order, before the `SendBill` Activity is
   called.
2. Immediately after the `UpdateInventory` block, add the compensating Activity,
`RevertInventory`, for this function. This will use a Go `defer` block:

```go
	defer func() {
		if err != nil {
			errCompensation := workflow.ExecuteActivity(ctx, RevertInventory, order.Items).Get(ctx, nil)
			err = multierr.Append(err, errCompensation)
		}
	}()
```

   This will ensure that, if the Workflow later encounters an error, it will be
   able to gracefully roll back the corresponding Activity. Using the `multierr`
   package lets you append errors incrementally in each `defer` block, so the
   Workflow will only run compensating Activities that correspond to the
   Activities that have already succceeded up to that point.
3. Next, add another `defer` block after `SendBill` to run the corresponding
   rollback Activity, `RefundCustomer`.

## Part C: Test the Rollback of Your Activities

Because this is fundamentally the same Exercise as #1 and #2, you can run the Worker and Starter the same way.

1. In one terminal, start the Worker by running:
   ```bash
   go run worker/main.go
   ```
2. In another terminal, start the Workflow by executing `start/main.go`. To
   trigger the error handling and rollback, run it with the bad credit card
   number again (otherwise it will just complete successfully, without
   demonstrating rollback):
   ```bash
   go run start/main.go --creditcard 1234
   ```
3. In your Web UI, you should see a `WorkflowExecutionFailed` Event to indicate
   that the Workflow failed. If you review the Workflow Timeline view, you'll
   see that after `SendBill` Activity succeeded, and the `ProcessCreditCard`
   Activity deliberately failed, we ran `RefundCustomer` and `RevertInventory`.

### This is the end of the exercise.
