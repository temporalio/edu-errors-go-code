## Exercise #2: Modifying Activity Options Using Non-Retryable Error Types and Heartbeats

During this exercise, you will:

- Configure non-retryable error types for Activities
- Implement customize retry policies for Activities
- Add Heartbeats and Heartbeat timeouts to help users monitor the health of Activities

Make your changes to the code in the `practice` subdirectory (look for
`TODO` comments that will guide you to where you should make changes to
the code). If you need a hint or want to verify your changes, look at
the complete version in the `solution` subdirectory.

## Part A: Convert Non-Retryable Errors to Be Handled By a Retry Policy

In this part of the exercise, you will modify the `ApplicationFailure` you defined
in `ProcessCreditCard` method in the first exercise to not be set as non-retryable
by default. After consideration, you've determined that while you may want to
immediately fail your Workflow on failure, others who call your Activity may not.

1. Open `activities.go`. In the `ProcessCreditCard` method from the last
exercise, modify the exception that is being thrown to be retryable. To do
   this, change `NewNonRetryableApplicationError` to `NewApplicationError`. Now,
   when an error is thrown from this Activity, the Activity will be retried.
2. Save and close the file.
3. Verify that your Error is now being retried by attempting to execute the
   Workflow. In one terminal, start the Worker by running:
   ```bash
   go run worker/main.go
   ```
4. In another terminal, start the Workflow by executing `start/main.go` with bad
   data to cause the error:
   ```bash
   go run start/main.go --creditcard 1234
   ```
5. Go to the WebUI and view the status of the Workflow. It should be
   **Running**. Inspect the Workflow and see that it is currently retrying the
   exception, verifying that the exception is no longer non-retryable. You can
   terminate this Workflow in the WebUI, as it will never successfully complete.
6. Stop your Worker by using **Ctrl-C** in the terminal it is running in.

## Part B: Configure Retry Policies to set Non-Retryable Error Types

Now that the exception from the `ProcessCreditCard` Activity is no longer set to
non-retryable, anyone who is running your Activity code may decide how to handle
the failure. In this case, imagine you have decided that you do not want the
Activity to retry upon failure -- but you don't want to "hardcode" it to be
non-retryable as in Exercise 1. In this part of the exercise, you will configure
a Retry Policy to disallow this using non-retryable error types.

Recall that a Retry Policy has the following attributes:

- Initial Interval: Amount of time that must elapse before the first retry occurs
- Backoff Coefficient: How much the retry interval increases (default is 2.0)
- Maximum Interval: The maximum interval between retries
- Maximum Attempts: The maximum number of execution attempts that can be made in the presence of failures

You can also specify errors types that are not retryable in the Retry Policy.
These are known as non-retryable error types. In the `ProcessCreditCard`
Activity in `activities.go`, you returned an `ApplicationFailure` with the
type `CreditCardError`. Now you will specify that error type as
non-retryable.

1. Open `workflow.go`. A `RetryPolicy` has already been defined with the
   following configuration:

```go
	retrypolicy := &temporal.RetryPolicy{
		MaximumInterval: time.Second * 10,
		MaximumAttempts: 3,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 120,
		RetryPolicy:         retrypolicy,
	}
```

2. Add an additional parameter to `retrypolicy` like so:

```go
NonRetryableErrorTypes: []string{"CreditCardError"},
```

3. Save and close the file.
4. Verify that your Error is once again failing the Workflow. In one terminal,
   start the Worker by running:
   ```bash
   go run worker/main.go
   ```
5. In another terminal, start the Workflow by executing `start/main.go` with bad
   data to cause the error:
   ```bash
   go run start/main.go --creditcard 1234
   ```
6. Go to the WebUI and view the status of the Workflow. You should see an
   `ActivityTaskFailed` error.

## Part C: Add Heartbeats

In this part of the exercise, you will add heartbeating to a new Activity,
`NotifyDeliveryDriver` The `NotifyDeliveryDriver` method attempts to contact a
driver to deliver the customers pizza. It may take a while for a delivery driver
to accept the delivery, and you want to ensure that the Activity is still alive
and processing. Heartbeats are used to do this, and fail fast if a failure is
detected.

In this exercise, instead of attempting to call an external service, you will
simulate a successful call to the `NotifyDeliveryDriver` method.

**How the simulation works**: The simulation starts by generating a number from
0 - 14. From there a loop is iterated over from 0 < 10, each time checking to
see if the random number matches the loop counter and then sleeping for 5 seconds.
Each iteration of the loop sends a heartbeat back letting the Workflow know that
progress is still being made. If the number matches a loop counter, it is a success
and `true` is returned. If it doesn't, then a delivery driver was unable to be
contacted and false is returned and the `status` of the `OrderConfirmation` will
be updated to reflect this.

1. Open `workflow.go`. Locate the `NotifyDeliveryDriver` function and uncomment
   it.
2. Save and close the file.
3. Next, open `activities.go`. Within the `NotifyDeliveryDriver` loop, above the
   `logger` call, add a heartbeat, providing the iteration number as the
   details.
   ```go
      activity.RecordHeartbeat(ctx, x)
   ```
4. Save and close the file.

## Part D: Add a Heartbeat Timeout

In the previous part of the exercise, you added a Heartbeat to an Activity. However,
you didn't set how long the Heartbeat should be inactive for before it is considered
a failed Heartbeat.

In this section, we will add a Heartbeat Timeout to your Activities.

1. Open `workflow.go`.
2. In your `ActivityOptions`, set the HeartbeatTimeout to a duration of 10
   seconds with `HeartbeatTimeout: 10 * time.Second,`. This sets the maximum
   time allowed between Activity Heartbeats before the Heartbeat is considered
   failed.
3. Save and close the file.

## Part E: Run the Workflow

Now you will run the Workflow and witness the Heartbeats happening in the
Web UI.

1. In one terminal, start the Worker by running:
   ```bash
   go run worker/main.go
   ```
2. In another terminal, start the Workflow by executing `start/main.go`:
   ```bash
   go run start/main.go
   ```
3. Now, go to the WebUI and find your workflow, which should be in the `Running`
   state. Click on it to enter the details page. Once you see `Heartbeat: <A_NUMBER>` 
   in your Worker output refresh the page and look for a **Pending Activities** 
   section. In this section you should see **Heartbeat Details** and JSON representing
   the payload. Remember, the simulation will finish at a random interval. You may
   need to run this a few times to see the results.

You have now seen how heartbeats are implemented and appear when an Activity is
running.

## (Optional) Part F: Failing a Heartbeat

Now that you've seen what a successful Heartbeat looks like, you should experience
a Heartbeat that is timing out. 

1. Open `activities.go`.
2. In the `NotifyDeliveryDriver` method, update the duration in the
   `workflow.Sleep()` call from 5s to 15s. This is longer than the Heartbeat
   Timeout you set in Step D.
3. Stop the Worker from the previous exercise and restart it:
   ```bash
   go run worker/main.go
   ```
4. In another terminal, start the Workflow by executing `start/main.go`:
   ```bash
   go run start/main.go
   ```
5. Once you see the first Heartbeat message appear in the logs, wait 15 seconds
   and refresh the WebUI. You should see the same **Pending Activities**
   section, but now there is a failure indicating that the Heartbeat timed out.
   You will also see how many retries are left and how long until the next
   retry. If the Activity isn't fixed before the final attempt it will fail.

You have now seen what happens when a Heartbeat times out.

### This is the end of the exercise.
