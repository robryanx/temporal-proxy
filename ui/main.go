package main

import (
	"context"
	"fmt"
	"log"
	"time"

	synchronousproxy "temporal-proxy"
	"temporal-proxy/proxy"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	fmt.Println("T-Shirt Order")

	status, err := CreateOrder(c)
	if err != nil {
		log.Fatalln("Unable to create order", err)
	}

	email := "test"
	for i:=0; i<10; i++ {
		status, err = UpdateOrder(c, status.OrderID, synchronousproxy.RegisterStage, email)
		if err != nil {
			log.Println("invalid email", err)
			continue
		}
	}

	fmt.Println("Thanks for your order!")
	fmt.Println("You will receive an email with shipping details shortly")
}

func CreateOrder(c client.Client) (synchronousproxy.OrderStatus, error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "ui-driven",
	}
	ctx := context.Background()
	var status synchronousproxy.OrderStatus

	we, err := c.ExecuteWorkflow(ctx, workflowOptions, synchronousproxy.OrderWorkflow)
	if err != nil {
		return status, fmt.Errorf("unable to execute order workflow: %w", err)
	}

	status.OrderID = we.GetID()

	return status, nil
}

func UpdateOrder(c client.Client, orderID string, stage string, value string) (synchronousproxy.OrderStatus, error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "ui-driven",
	}
	ctx := context.Background()
	status := synchronousproxy.OrderStatus{OrderID: orderID}

	
	we, err := c.ExecuteWorkflow(ctx, workflowOptions, synchronousproxy.UpdateOrderWorkflow, orderID, stage, value)
	if err != nil {
		return status, fmt.Errorf("unable to execute workflow: %w", err)
	}

	time.Sleep(time.Second)

	start := time.Now()
	proxy.SendExternalRequest(c, context.Background(), we.GetID(), "test", "test")
	err = we.Get(ctx, &status)
	if err != nil {
		return status, err
	}
	fmt.Println(time.Now().Sub(start))

	return status, nil
}
