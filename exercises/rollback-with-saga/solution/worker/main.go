package main

import (
	pizza "errstrat/exercises/rollback-with-saga/solution"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, pizza.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(pizza.PizzaWorkflow)
	w.RegisterActivity(pizza.GetDistance)
	w.RegisterActivity(pizza.SendBill)
	w.RegisterActivity(pizza.ProcessCreditCard)
	w.RegisterActivity(pizza.UpdateInventory)
	w.RegisterActivity(pizza.RevertInventory)
	w.RegisterActivity(pizza.RefundCustomer)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}
