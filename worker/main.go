package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	synchronousproxy "temporal-proxy"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "ui-driven", worker.Options{})

	w.RegisterWorkflow(synchronousproxy.OrderWorkflow)
	w.RegisterWorkflow(synchronousproxy.UpdateOrderWorkflow)
	w.RegisterActivity(synchronousproxy.RegisterEmail)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
